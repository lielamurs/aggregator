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

	for result := range results {
		if result.Err != nil {
			logger.WithError(result.Err).WithField("bank", result.BankName).Error("Bank submission failed")

			if err := s.saveBankSubmission(ctx, customerApp.ID, result.BankName, dto.SubmissionStatusFailed, "", result.Err); err != nil {
				logger.WithError(err).WithField("bank", result.BankName).Error("Failed to save bank submission")
			}
		} else {
			logger.WithField("bank", result.BankName).Info("Bank submission successful")

			if err := s.saveBankSubmission(ctx, customerApp.ID, result.BankName, dto.SubmissionStatusDraft, result.SubmissionID, nil); err != nil {
				logger.WithError(err).WithField("bank", result.BankName).Error("Failed to save bank submission")
			}
		}
	}

	logger.Info("Application processing completed - bank submissions saved as DRAFT")
}

func (s *applicationService) submitToBankAsync(ctx context.Context, bank BankService, customerApp *dto.CustomerApplication, results chan<- dto.BankResult) {
	logger := s.logger.WithFields(logrus.Fields{
		"application_id": customerApp.ID,
		"bank":           bank.GetBankName(),
	})

	logger.Info("Submitting application to bank")

	response, err := bank.SubmitApplication(ctx, customerApp.CustomerData)
	if err != nil {
		logger.WithError(err).Error("Bank submission failed")
		results <- dto.BankResult{
			BankName: bank.GetBankName(),
			Err:      err,
		}
		return
	}

	logger.WithField("submission_id", response.ID).Info("Bank submission successful")
	results <- dto.BankResult{
		BankName:     bank.GetBankName(),
		SubmissionID: response.ID,
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

func (s *applicationService) saveBankSubmission(ctx context.Context, applicationID uuid.UUID, bankName string, status dto.BankSubmissionStatus, bankID string, submissionErr error) error {
	exists, err := s.applicationsRepo.Exists(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to verify application exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("application with ID %s not found", applicationID)
	}

	now := time.Now()
	bankSubmission := &dto.BankSubmission{
		ID:        uuid.New(),
		BankName:  bankName,
		Status:    status,
		BankID:    bankID,
		CreatedAt: now,
	}

	if status == dto.SubmissionStatusDraft {
		bankSubmission.SubmittedAt = now
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
