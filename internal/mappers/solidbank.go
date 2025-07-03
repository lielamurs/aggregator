package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
)

func ToSolidBankRequestFromApplicationRequest(req dto.ApplicationRequest) *dto.SolidBankApplicationRequest {
	return &dto.SolidBankApplicationRequest{
		Phone:           req.Phone,
		Email:           req.Email,
		MonthlyIncome:   req.MonthlyIncome,
		MonthlyExpenses: req.MonthlyExpenses,
		MaritalStatus:   req.MaritalStatus,
		AgreeToBeScored: req.AgreeToBeScored,
		Amount:          req.Amount,
	}
}

func ToOfferFromSolidBankApplication(app dto.SolidBankApplication, bankName string) *dto.Offer {
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
