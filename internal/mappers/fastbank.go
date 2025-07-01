package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
)

func ToFastBankRequest(ar *dto.ApplicationRequest) dto.FastBankApplicationRequest {
	if ar == nil {
		return dto.FastBankApplicationRequest{}
	}

	return dto.FastBankApplicationRequest{
		PhoneNumber:              ar.Phone,
		Email:                    ar.Email,
		MonthlyIncomeAmount:      ar.MonthlyIncome,
		MonthlyCreditLiabilities: ar.MonthlyExpenses,
		Dependents:               ar.Dependents,
		AgreeToDataSharing:       ar.AgreeToBeScored,
		Amount:                   ar.Amount,
	}
}

func FastBankOfferToOffer(fbo *dto.FastBankOffer, bankName string) dto.Offer {
	if fbo == nil {
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
		MonthlyPaymentAmount: fbo.MonthlyPaymentAmount,
		TotalRepaymentAmount: fbo.TotalRepaymentAmount,
		NumberOfPayments:     fbo.NumberOfPayments,
		AnnualPercentageRate: fbo.AnnualPercentageRate,
		FirstRepaymentDate:   fbo.FirstRepaymentDate,
		Status:               dto.OfferStatusApproved,
		CreatedAt:            time.Now(),
	}
}

func FastBankApplicationToOffers(app *dto.FastBankApplication, bankName string) []dto.Offer {
	var offers []dto.Offer

	if app != nil && app.Offer != nil {
		offer := FastBankOfferToOffer(app.Offer, bankName)
		offers = append(offers, offer)
	}

	return offers
}
