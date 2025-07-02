package models

import (
	"time"

	"github.com/google/uuid"
)

type Offer struct {
	ID                   uuid.UUID
	ApplicationID        uuid.UUID
	BankName             string
	MonthlyPaymentAmount *float64
	TotalRepaymentAmount *float64
	NumberOfPayments     *int
	AnnualPercentageRate *float64
	FirstRepaymentDate   *string
	Status               string
	CreatedAt            time.Time
}
