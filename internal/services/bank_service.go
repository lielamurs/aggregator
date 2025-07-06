package services

import (
	"context"

	"github.com/lielamurs/aggregator/internal/dto"
)

type BankService interface {
	GetBankName() string
	SubmitApplication(ctx context.Context, req dto.ApplicationRequest) (*dto.BankSubmissionResponse, error)
	GetOffer(ctx context.Context, bankID string) (*dto.Offer, error)
}
