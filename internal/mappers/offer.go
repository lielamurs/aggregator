package mappers

import (
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/models"
)

func ToOfferModel(offer *dto.Offer) *models.Offer {
	if offer == nil {
		return nil
	}

	return &models.Offer{
		ID:                   offer.ID,
		BankName:             offer.BankName,
		MonthlyPaymentAmount: offer.MonthlyPaymentAmount,
		TotalRepaymentAmount: offer.TotalRepaymentAmount,
		NumberOfPayments:     offer.NumberOfPayments,
		AnnualPercentageRate: offer.AnnualPercentageRate,
		FirstRepaymentDate:   offer.FirstRepaymentDate,
		Status:               string(offer.Status),
		CreatedAt:            offer.CreatedAt,
	}
}

func ToOfferFromModel(offer *models.Offer) *dto.Offer {
	if offer == nil {
		return nil
	}

	return &dto.Offer{
		ID:                   offer.ID,
		BankName:             offer.BankName,
		MonthlyPaymentAmount: offer.MonthlyPaymentAmount,
		TotalRepaymentAmount: offer.TotalRepaymentAmount,
		NumberOfPayments:     offer.NumberOfPayments,
		AnnualPercentageRate: offer.AnnualPercentageRate,
		FirstRepaymentDate:   offer.FirstRepaymentDate,
		Status:               dto.OfferStatus(offer.Status),
		CreatedAt:            offer.CreatedAt,
	}
}
