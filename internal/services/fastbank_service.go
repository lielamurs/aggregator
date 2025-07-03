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

type fastBankService struct {
	config     config.FastBankConfig
	httpClient *HTTPClient
	logger     *logrus.Logger
}

func NewFastBankService(config config.FastBankConfig, logger *logrus.Logger) BankService {
	return &fastBankService{
		config:     config,
		httpClient: NewHTTPClient(time.Duration(config.Timeout)*time.Second, logger),
		logger:     logger,
	}
}

func (s *fastBankService) GetBankName() string {
	return "FastBank"
}

func (s *fastBankService) SubmitApplication(ctx context.Context, req dto.ApplicationRequest) (*dto.Offer, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"bank":   "FastBank",
		"phone":  req.Phone,
		"amount": req.Amount,
	})

	fastBankReq := mappers.ToFastBankRequestFromApplicationRequest(req)
	submitURL := fmt.Sprintf("%s/applications", s.config.BaseURL)
	var fastBankApp dto.FastBankApplication

	err := s.httpClient.PostJSON(ctx, submitURL, fastBankReq, &fastBankApp)
	if err != nil {
		logger.WithError(err).Error("Failed to submit application to FastBank")
		return nil, fmt.Errorf("FastBank submission failed: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"application_id": fastBankApp.ID,
		"status":         fastBankApp.Status,
	}).Info("FastBank application submitted")

	if fastBankApp.Status == "PROCESSED" {
		return s.handleProcessedApplication(fastBankApp, logger, s.GetBankName())
	}

	pollURL := fmt.Sprintf("%s/applications/%s", s.config.BaseURL, fastBankApp.ID)
	offer, err := s.pollForCompletion(ctx, pollURL, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to poll FastBank application")
		return nil, fmt.Errorf("FastBank polling failed: %w", err)
	}

	logger.WithField("offer_id", offer.ID).Info("FastBank application completed successfully")
	return offer, nil
}

func (s *fastBankService) handleProcessedApplication(fastBankApp dto.FastBankApplication, logger *logrus.Entry, bankName string) (*dto.Offer, error) {
	if fastBankApp.Offer != nil {
		logger.Info("FastBank application already processed with offer")
	} else {
		logger.Info("FastBank application already processed but rejected")
	}
	offer := mappers.ToOfferFromFastBankApplication(fastBankApp, bankName)
	if offer == nil {
		logger.Error("FastBank application mapping returned nil offer")
		return nil, fmt.Errorf("failed to map FastBank application")
	}
	return offer, nil
}

func (s *fastBankService) pollForCompletion(ctx context.Context, url string, logger *logrus.Entry) (*dto.Offer, error) {
	maxAttempts := s.config.MaxAttempts
	pollInterval := time.Duration(s.config.PollInterval) * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.WithField("attempt", attempt).Debug("Polling FastBank application status")

		var fastBankApp dto.FastBankApplication
		err := s.httpClient.GetJSON(ctx, url, &fastBankApp)
		if err != nil {
			logger.WithError(err).WithField("attempt", attempt).Error("Failed to poll FastBank application")
			return nil, fmt.Errorf("polling failed: %w", err)
		}

		if fastBankApp.Status == "PROCESSED" {
			return s.handleProcessedApplication(fastBankApp, logger, s.GetBankName())
		}

		if attempt < maxAttempts {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled while polling: %w", ctx.Err())
			case <-time.After(pollInterval):
			}
		}
	}

	return nil, fmt.Errorf("FastBank application did not complete within timeout period")
}
