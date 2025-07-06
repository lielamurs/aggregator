package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/mappers"
	"github.com/lielamurs/aggregator/internal/models"
	"github.com/lielamurs/aggregator/internal/repository"
	"github.com/sirupsen/logrus"
)

type SubmissionService interface {
	ProcessSubmissions(ctx context.Context) error
}

type submissionService struct {
	applicationsRepo    *repository.ApplicationsRepository
	offersRepo          *repository.OffersRepository
	bankSubmissionsRepo *repository.BankSubmissionsRepository
	bankServices        []BankService
	logger              *logrus.Logger
}

func NewSubmissionService(
	applicationsRepo *repository.ApplicationsRepository,
	offersRepo *repository.OffersRepository,
	bankSubmissionsRepo *repository.BankSubmissionsRepository,
	bankServices []BankService,
	logger *logrus.Logger,
) SubmissionService {
	return &submissionService{
		applicationsRepo:    applicationsRepo,
		offersRepo:          offersRepo,
		bankSubmissionsRepo: bankSubmissionsRepo,
		bankServices:        bankServices,
		logger:              logger,
	}
}

func (s *submissionService) ProcessSubmissions(ctx context.Context) error {
	logger := s.logger.WithField("component", "submission_processor")
	logger.Info("Starting submission processing cycle")

	processingApplications, err := s.applicationsRepo.GetProcessingApplicationsWithBankSubmissions(ctx)
	if err != nil {
		logger.WithError(err).Error("Failed to get processing applications")
		return fmt.Errorf("failed to get processing applications: %w", err)
	}

	logger.WithField("count", len(processingApplications)).Info("Found processing applications")

	for _, app := range processingApplications {
		if err := s.processApplication(ctx, &app); err != nil {
			logger.WithError(err).WithField("application_id", app.ID).Error("Failed to process application")
		}
	}

	logger.Info("Submission processing cycle completed")
	return nil
}

func (s *submissionService) processApplication(ctx context.Context, app *models.Application) error {
	logger := s.logger.WithField("application_id", app.ID)
	logger.Info("Processing application submissions")

	draftSubmissions := make([]models.BankSubmission, 0)
	for _, submission := range app.BankSubmissions {
		if submission.Status == string(dto.SubmissionStatusDraft) {
			draftSubmissions = append(draftSubmissions, submission)
		}
	}

	if len(draftSubmissions) == 0 {
		logger.Info("No draft submissions found for application")
		return nil
	}

	logger.WithField("draft_count", len(draftSubmissions)).Info("Found draft submissions")

	allCompleted := true
	for _, submission := range draftSubmissions {
		if err := s.processSubmission(ctx, app.ID, &submission); err != nil {
			logger.WithError(err).WithField("bank", submission.BankName).Error("Failed to process submission")
			allCompleted = false
		}
	}

	if allCompleted {
		if err := s.updateApplicationStatus(ctx, app.ID, dto.StatusCompleted); err != nil {
			logger.WithError(err).Error("Failed to update application status to completed")
			return fmt.Errorf("failed to update application status: %w", err)
		}
		logger.Info("Application processing completed")
	}

	return nil
}

func (s *submissionService) processSubmission(ctx context.Context, applicationID uuid.UUID, submission *models.BankSubmission) error {
	logger := s.logger.WithFields(logrus.Fields{
		"application_id": applicationID,
		"bank":           submission.BankName,
		"submission_id":  submission.ID,
		"bank_id":        submission.BankID,
	})

	var bankService BankService
	for _, service := range s.bankServices {
		if service.GetBankName() == submission.BankName {
			bankService = service
			break
		}
	}

	if bankService == nil {
		logger.Error("Bank service not found")
		return fmt.Errorf("bank service not found for %s", submission.BankName)
	}

	if submission.BankID == nil {
		logger.Error("Bank ID is nil, cannot get offer")
		return fmt.Errorf("bank ID is nil for submission %s", submission.ID)
	}

	offer, err := bankService.GetOffer(ctx, *submission.BankID)
	if err != nil {
		logger.WithError(err).Error("Failed to get offer from bank")

		submission.Status = string(dto.SubmissionStatusFailed)
		errorMsg := err.Error()
		submission.ErrorMessage = &errorMsg
		now := time.Now()
		submission.CompletedAt = &now

		if updateErr := s.bankSubmissionsRepo.Update(ctx, submission); updateErr != nil {
			logger.WithError(updateErr).Error("Failed to update failed submission")
		}

		return fmt.Errorf("failed to get offer: %w", err)
	}

	if offer == nil {
		logger.Debug("Application not yet processed by bank")
		return nil
	}

	logger.Info("Successfully retrieved offer from bank")

	if err := s.saveOffer(ctx, applicationID, offer); err != nil {
		logger.WithError(err).Error("Failed to save offer")
		return fmt.Errorf("failed to save offer: %w", err)
	}

	submission.Status = string(dto.SubmissionStatusSuccess)
	now := time.Now()
	submission.CompletedAt = &now

	if err := s.bankSubmissionsRepo.Update(ctx, submission); err != nil {
		logger.WithError(err).Error("Failed to update successful submission")
		return fmt.Errorf("failed to update submission: %w", err)
	}

	logger.Info("Submission processed successfully")
	return nil
}

func (s *submissionService) saveOffer(ctx context.Context, applicationID uuid.UUID, bankOffer *dto.Offer) error {
	if bankOffer == nil {
		return fmt.Errorf("offer cannot be nil")
	}

	bankOffer.ID = uuid.New()

	offer := mappers.ToOfferModel(bankOffer)
	if offer == nil {
		return fmt.Errorf("failed to convert offer to model")
	}

	offer.ApplicationID = applicationID
	return s.offersRepo.Create(ctx, offer)
}

func (s *submissionService) updateApplicationStatus(ctx context.Context, applicationID uuid.UUID, status dto.ApplicationStatus) error {
	app, err := s.applicationsRepo.GetByID(ctx, applicationID)
	if err != nil {
		return fmt.Errorf("failed to get application: %w", err)
	}

	app.Status = string(status)
	app.UpdatedAt = time.Now()

	return s.applicationsRepo.Update(ctx, app)
}
