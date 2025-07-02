package models

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID              uuid.UUID
	Phone           string
	Email           string
	MonthlyIncome   float64
	MonthlyExpenses float64
	MaritalStatus   string
	AgreeToBeScored bool
	Amount          float64
	Dependents      int
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Offers          []Offer          `gorm:"foreignKey:ApplicationID"`
	BankSubmissions []BankSubmission `gorm:"foreignKey:ApplicationID"`
}
