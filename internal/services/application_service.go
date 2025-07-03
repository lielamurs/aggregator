package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/mappers"
	"github.com/lielamurs/aggregator/internal/models"
	"github.com/lielamurs/aggregator/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ApplicationService interface {
	SubmitApplication(ctx context.Context, app *dto.CustomerApplication) (*dto.ApplicationResponse, error)
	GetApplicationStatus(ctx context.Context, applicationID uuid.UUID) (*models.Application, error)
}

type applicationService struct {
	applicationsRepo    *repository.ApplicationsRepository
	offersRepo          *repository.OffersRepository
	bankSubmissionsRepo *repository.BankSubmissionsRepository
	bankServices        []BankService
	logger              *logrus.Logger
}

func NewApplicationService(
	applicationsRepo *repository.ApplicationsRepository,
	offersRepo *repository.OffersRepository,
	bankSubmissionsRepo *repository.BankSubmissionsRepository,
	bankServices []BankService,
	logger *logrus.Logger,
) ApplicationService {
	return &applicationService{
		applicationsRepo:    applicationsRepo,
		offersRepo:          offersRepo,
		bankSubmissionsRepo: bankSubmissionsRepo,
		bankServices:        bankServices,
		logger:              logger,
	}
}

func (s *applicationService) SubmitApplication(ctx context.Context, customerApp *dto.CustomerApplication) (*dto.ApplicationResponse, error) {
	application := mappers.ToApplicationModel(customerApp)
	if application == nil {
		return nil, fmt.Errorf("failed to convert application to model")
	}

	if err := s.applicationsRepo.Create(ctx, application); err != nil {
		s.logger.WithError(err).WithField("application_id", customerApp.ID).Error("Failed to save application")
		return nil, fmt.Errorf("failed to save application: %w", err)
	}

	go s.processApplication(context.Background(), customerApp)

	return &dto.ApplicationResponse{
		ID:     customerApp.ID,
		Status: customerApp.Status,
	}, nil
}

func (s *applicationService) GetApplicationStatus(ctx context.Context, applicationID uuid.UUID) (*models.Application, error) {
	logger := s.logger.WithField("application_id", applicationID)

	application, err := s.applicationsRepo.GetByID(ctx, applicationID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Debug("Application not found")
			return nil, fmt.Errorf("application with ID %s not found", applicationID)
		}
		logger.WithError(err).Error("Failed to get application from database")
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return application, nil
}

func (s *applicationService) processApplication(ctx context.Context, customerApp *dto.CustomerApplication) {
	logger := s.logger.WithField("application_id", customerApp.ID)
	logger.Info("Starting application processing")

	customerApp.Status = dto.StatusProcessing
	customerApp.UpdatedAt = time.Now()

	if err := s.updateApplication(ctx, customerApp); err != nil {
		logger.WithError(err).Error("Failed to update application status to processing")
		return
	}

	var wg sync.WaitGroup
	results := make(chan dto.BankResult, len(s.bankServices))

	for _, bankService := range s.bankServices {
		wg.Add(1)
		go func(bank BankService) {
			defer wg.Done()
			s.submitToBankAsync(ctx, bank, customerApp, results)
		}(bankService)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	allSuccessful := true
	for result := range results {
		if result.Err != nil {
			logger.WithError(result.Err).WithField("bank", result.BankName).Error("Bank submission failed")
			allSuccessful = false

			if err := s.saveBankSubmission(ctx, customerApp.ID, result.BankName, dto.SubmissionStatusFailed, result.Err); err != nil {
				logger.WithError(err).WithField("bank", result.BankName).Error("Failed to save bank submission")
			}
		} else {
			logger.WithField("bank", result.BankName).Info("Bank submission successful")

			if err := s.saveBankSubmission(ctx, customerApp.ID, result.BankName, dto.SubmissionStatusSuccess, nil); err != nil {
				logger.WithError(err).WithField("bank", result.BankName).Error("Failed to save bank submission")
			}

			if err := s.saveOffer(ctx, customerApp.ID, result.Offer); err != nil {
				logger.WithError(err).WithField("bank", result.BankName).Error("Failed to save offer")
			}
		}
	}

	customerApp.Status = dto.StatusCompleted
	if allSuccessful {
		logger.Info("Application processing completed successfully")
	} else {
		logger.Info("Application processing completed with some bank failures")
	}

	customerApp.UpdatedAt = time.Now()
	if err := s.updateApplication(ctx, customerApp); err != nil {
		logger.WithError(err).Error("Failed to update final application status")
	}
}

func (s *applicationService) submitToBankAsync(ctx context.Context, bank BankService, customerApp *dto.CustomerApplication, results chan<- dto.BankResult) {
	logger := s.logger.WithFields(logrus.Fields{
		"application_id": customerApp.ID,
		"bank":           bank.GetBankName(),
	})

	logger.Info("Submitting application to bank")

	offer, err := bank.SubmitApplication(ctx, customerApp.CustomerData)
	if err != nil {
		logger.WithError(err).Error("Bank submission failed")
		results <- dto.BankResult{
			BankName: bank.GetBankName(),
			Err:      err,
		}
		return
	}

	logger.Info("Bank submission successful")
	results <- dto.BankResult{
		BankName: bank.GetBankName(),
		Offer:    offer,
	}
}

func (s *applicationService) updateApplication(ctx context.Context, customerApp *dto.CustomerApplication) error {
	if customerApp == nil {
		return fmt.Errorf("application cannot be nil")
	}

	application := mappers.ToApplicationModel(customerApp)
	if application == nil {
		return fmt.Errorf("failed to convert application to model")
	}

	return s.applicationsRepo.Update(ctx, application)
}

func (s *applicationService) saveOffer(ctx context.Context, applicationID uuid.UUID, bankOffer *dto.Offer) error {
	if bankOffer == nil {
		return fmt.Errorf("offer cannot be nil")
	}

	exists, err := s.applicationsRepo.Exists(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to verify application exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("application with ID %s not found", applicationID)
	}

	bankOffer.ID = uuid.New()

	offer := mappers.ToOfferModel(bankOffer)
	if offer == nil {
		return fmt.Errorf("failed to convert offer to model")
	}

	offer.ApplicationID = applicationID
	return s.offersRepo.Create(ctx, offer)
}

func (s *applicationService) saveBankSubmission(ctx context.Context, applicationID uuid.UUID, bankName string, status dto.BankSubmissionStatus, submissionErr error) error {
	exists, err := s.applicationsRepo.Exists(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to verify application exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("application with ID %s not found", applicationID)
	}

	bankSubmission := &dto.BankSubmission{
		ID:        uuid.New(),
		BankName:  bankName,
		Status:    status,
		CreatedAt: time.Now(),
	}

	if submissionErr != nil {
		errorMsg := submissionErr.Error()
		bankSubmission.ErrorMessage = &errorMsg
	}

	submission := mappers.ToBankSubmissionModel(bankSubmission)
	if submission == nil {
		return fmt.Errorf("failed to convert bank submission to model")
	}

	submission.ApplicationID = applicationID
	return s.bankSubmissionsRepo.Create(ctx, submission)
}
