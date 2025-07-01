package dto

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationRequest struct {
	Phone           string  `json:"phone" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	MonthlyIncome   float64 `json:"monthlyIncome" validate:"required,min=0"`
	MonthlyExpenses float64 `json:"monthlyExpenses" validate:"required,min=0"`
	MaritalStatus   string  `json:"maritalStatus" validate:"required,oneof=SINGLE MARRIED DIVORCED WIDOWED COHABITING"`
	AgreeToBeScored bool    `json:"agreeToBeScored" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,min=0"`
	Dependents      int     `json:"dependents" validate:"min=0"`
}

type Application struct {
	ID              uuid.UUID          `json:"id"`
	CustomerData    ApplicationRequest `json:"customerData"`
	Status          ApplicationStatus  `json:"status"`
	Offers          []Offer            `json:"offers"`
	BankSubmissions []BankSubmission   `json:"bankSubmissions"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

type ApplicationStatus string

const (
	StatusPending    ApplicationStatus = "PENDING"
	StatusProcessing ApplicationStatus = "PROCESSING"
	StatusCompleted  ApplicationStatus = "COMPLETED"
	StatusFailed     ApplicationStatus = "FAILED"
)

type ApplicationResponse struct {
	ID     uuid.UUID         `json:"id"`
	Status ApplicationStatus `json:"status"`
}

type ApplicationStatusResponse struct {
	ID              uuid.UUID         `json:"id"`
	Status          ApplicationStatus `json:"status"`
	Offers          []Offer           `json:"offers"`
	BankSubmissions []BankSubmission  `json:"bankSubmissions"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}
