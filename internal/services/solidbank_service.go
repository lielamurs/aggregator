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

func (s *solidBankService) SubmitApplication(ctx context.Context, req dto.ApplicationRequest) (*dto.BankSubmissionResponse, error) {
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

	return &dto.BankSubmissionResponse{
		ID:     solidBankApp.ID,
		Status: solidBankApp.Status,
	}, nil
}

func (s *solidBankService) GetOffer(ctx context.Context, bankID string) (*dto.Offer, error) {
	logger := s.logger.WithFields(logrus.Fields{
		"bank":    "SolidBank",
		"bank_id": bankID,
	})

	pollURL := fmt.Sprintf("%s/applications/%s", s.config.BaseURL, bankID)
	var solidBankApp dto.SolidBankApplication

	err := s.httpClient.GetJSON(ctx, pollURL, &solidBankApp)
	if err != nil {
		logger.WithError(err).Error("Failed to get SolidBank application")
		return nil, fmt.Errorf("SolidBank get application failed: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"application_id": solidBankApp.ID,
		"status":         solidBankApp.Status,
	}).Info("SolidBank application status retrieved")

	if solidBankApp.Status == "PROCESSED" {
		return s.handleProcessedApplication(solidBankApp, logger, s.GetBankName())
	}

	return nil, nil
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
