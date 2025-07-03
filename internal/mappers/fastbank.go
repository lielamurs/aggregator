package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
)

func ToFastBankRequestFromApplicationRequest(req dto.ApplicationRequest) *dto.FastBankApplicationRequest {
	return &dto.FastBankApplicationRequest{
		PhoneNumber:              req.Phone,
		Email:                    req.Email,
		MonthlyIncomeAmount:      req.MonthlyIncome,
		MonthlyCreditLiabilities: req.MonthlyExpenses,
		Dependents:               req.Dependents,
		AgreeToDataSharing:       req.AgreeToBeScored,
		Amount:                   req.Amount,
	}
}

func ToOfferFromFastBankApplication(app dto.FastBankApplication, bankName string) *dto.Offer {
	if app.Status != "PROCESSED" {
		return nil
	}

	offer := &dto.Offer{
		ID:        uuid.New(),
		BankName:  bankName,
		CreatedAt: time.Now(),
	}

	if app.Offer != nil {
		offer.Status = dto.OfferStatusApproved
		offer.MonthlyPaymentAmount = &app.Offer.MonthlyPaymentAmount
		offer.TotalRepaymentAmount = &app.Offer.TotalRepaymentAmount
		offer.NumberOfPayments = &app.Offer.NumberOfPayments
		offer.AnnualPercentageRate = &app.Offer.AnnualPercentageRate
		offer.FirstRepaymentDate = &app.Offer.FirstRepaymentDate
	} else {
		offer.Status = dto.OfferStatusRejected
	}

	return offer
}
