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

func (s *fastBankService) SubmitApplication(ctx context.Context, req dto.ApplicationRequest) (*dto.BankSubmissionResponse, error) {
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

	return &dto.BankSubmissionResponse{
		ID:     fastBankApp.ID,
		Status: fastBankApp.Status,
	}, nil
}

func (s *fastBankService) GetOffer(ctx context.Context, bankID string) (*dto.Offer, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"bank":    "FastBank",
		"bank_id": bankID,
	})

	pollURL := fmt.Sprintf("%s/applications/%s", s.config.BaseURL, bankID)
	var fastBankApp dto.FastBankApplication

	err := s.httpClient.GetJSON(ctx, pollURL, &fastBankApp)
	if err != nil {
		logger.WithError(err).Error("Failed to get FastBank application")
		return nil, fmt.Errorf("FastBank get application failed: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"application_id": fastBankApp.ID,
		"status":         fastBankApp.Status,
	}).Info("FastBank application status retrieved")

	if fastBankApp.Status == "PROCESSED" {
		return s.handleProcessedApplication(fastBankApp, logger, s.GetBankName())
	}

	return nil, nil
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
