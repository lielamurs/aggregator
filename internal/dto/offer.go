package dto

import (
	"time"

	"github.com/google/uuid"
)

type Offer struct {
	ID                   uuid.UUID   `json:"id"`
	BankName             string      `json:"bankName"`
	MonthlyPaymentAmount *float64    `json:"monthlyPaymentAmount,omitempty"`
	TotalRepaymentAmount *float64    `json:"totalRepaymentAmount,omitempty"`
	NumberOfPayments     *int        `json:"numberOfPayments,omitempty"`
	AnnualPercentageRate *float64    `json:"annualPercentageRate,omitempty"`
	FirstRepaymentDate   *string     `json:"firstRepaymentDate,omitempty"`
	Status               OfferStatus `json:"status"`
	CreatedAt            time.Time   `json:"createdAt"`
}

type OfferStatus string

const (
	OfferStatusApproved OfferStatus = "APPROVED"
	OfferStatusRejected OfferStatus = "REJECTED"
)
