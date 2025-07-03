package services

import (
	"context"
	"fmt"
	"time"

	"github.com/lielamurs/aggregator/internal/config"
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/mappers"
	"github.com/sirupsen/logrus"
)

type solidBankService struct {
	config     config.SolidBankConfig
	httpClient *HTTPClient
	logger     *logrus.Logger
}

func NewSolidBankService(config config.SolidBankConfig, logger *logrus.Logger) BankService {
	return &solidBankService{
		config:     config,
		httpClient: NewHTTPClient(time.Duration(config.Timeout)*time.Second, logger),
		logger:     logger,
	}
}

func (s *solidBankService) GetBankName() string {
	return "SolidBank"
}

func (s *solidBankService) SubmitApplication(ctx context.Context, req dto.ApplicationRequest) (*dto.Offer, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"bank":   "SolidBank",
		"phone":  req.Phone,
		"amount": req.Amount,
	})

	solidBankReq := mappers.ToSolidBankRequestFromApplicationRequest(req)
	submitURL := fmt.Sprintf("%s/applications", s.config.BaseURL)
	var solidBankApp dto.SolidBankApplication

	err := s.httpClient.PostJSON(ctx, submitURL, solidBankReq, &solidBankApp)
	if err != nil {
		logger.WithError(err).Error("Failed to submit application to SolidBank")
		return nil, fmt.Errorf("SolidBank submission failed: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"application_id": solidBankApp.ID,
		"status":         solidBankApp.Status,
	}).Info("SolidBank application submitted")

	if solidBankApp.Status == "PROCESSED" {
		return s.handleProcessedApplication(solidBankApp, logger, s.GetBankName())
	}

	pollURL := fmt.Sprintf("%s/applications/%s", s.config.BaseURL, solidBankApp.ID)
	offer, err := s.pollForCompletion(ctx, pollURL, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to poll SolidBank application")
		return nil, fmt.Errorf("SolidBank polling failed: %w", err)
	}

	logger.WithField("offer_id", offer.ID).Info("SolidBank application completed successfully")
	return offer, nil
}

func (s *solidBankService) handleProcessedApplication(solidBankApp dto.SolidBankApplication, logger *logrus.Entry, bankName string) (*dto.Offer, error) {
	if solidBankApp.Offer != nil {
		logger.Info("SolidBank application already processed with offer")
	} else {
		logger.Info("SolidBank application already processed but rejected")
	}
	offer := mappers.ToOfferFromSolidBankApplication(solidBankApp, bankName)
	if offer == nil {
		logger.Error("SolidBank application mapping returned nil offer")
		return nil, fmt.Errorf("failed to map SolidBank application")
	}
	return offer, nil
}

func (s *solidBankService) pollForCompletion(ctx context.Context, url string, logger *logrus.Entry) (*dto.Offer, error) {
	maxAttempts := s.config.MaxAttempts
	pollInterval := time.Duration(s.config.PollInterval) * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.WithField("attempt", attempt).Debug("Polling SolidBank application status")

		var solidBankApp dto.SolidBankApplication
		err := s.httpClient.GetJSON(ctx, url, &solidBankApp)
		if err != nil {
			logger.WithError(err).WithField("attempt", attempt).Error("Failed to poll SolidBank application")
			return nil, fmt.Errorf("polling failed: %w", err)
		}

		logger.WithFields(logrus.Fields{
			"attempt":   attempt,
			"status":    solidBankApp.Status,
			"has_offer": solidBankApp.Offer != nil,
		}).Debug("SolidBank application status polled")

		if solidBankApp.Status == "PROCESSED" {
			return s.handleProcessedApplication(solidBankApp, logger, s.GetBankName())
		}

		if attempt < maxAttempts {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled while polling: %w", ctx.Err())
			case <-time.After(pollInterval):
			}
		}
	}

	return nil, fmt.Errorf("SolidBank application did not complete within timeout period")
}
