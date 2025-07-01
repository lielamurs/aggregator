package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
)

func ToSolidBankRequest(ar *dto.ApplicationRequest) dto.SolidBankApplicationRequest {
	if ar == nil {
		return dto.SolidBankApplicationRequest{}
	}

	return dto.SolidBankApplicationRequest{
		Phone:           ar.Phone,
		Email:           ar.Email,
		MonthlyIncome:   ar.MonthlyIncome,
		MonthlyExpenses: ar.MonthlyExpenses,
		MaritalStatus:   ar.MaritalStatus,
		AgreeToBeScored: ar.AgreeToBeScored,
		Amount:          ar.Amount,
	}
}

func SolidBankOfferToOffer(sbo *dto.SolidBankOffer, bankName string) dto.Offer {
	if sbo == nil {
		return dto.Offer{
			ID:        uuid.New(),
			BankName:  bankName,
			Status:    dto.OfferStatusRejected,
			CreatedAt: time.Now(),
		}
	}

	return dto.Offer{
		ID:                   uuid.New(),
		BankName:             bankName,
		MonthlyPaymentAmount: sbo.MonthlyPaymentAmount,
		TotalRepaymentAmount: sbo.TotalRepaymentAmount,
		NumberOfPayments:     sbo.NumberOfPayments,
		AnnualPercentageRate: sbo.AnnualPercentageRate,
		FirstRepaymentDate:   sbo.FirstRepaymentDate,
		Status:               dto.OfferStatusApproved,
		CreatedAt:            time.Now(),
	}
}

func SolidBankApplicationToOffers(app *dto.SolidBankApplication, bankName string) []dto.Offer {
	var offers []dto.Offer

	if app != nil && app.Offer != nil {
		offer := SolidBankOfferToOffer(app.Offer, bankName)
		offers = append(offers, offer)
	}

	return offers
}
