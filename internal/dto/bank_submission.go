package dto

import (
	"time"

	"github.com/google/uuid"
)

type BankSubmission struct {
	ID           uuid.UUID            `json:"id"`
	BankName     string               `json:"bankName"`
	Status       BankSubmissionStatus `json:"status"`
	BankID       string               `json:"bankId,omitempty"`
	SubmittedAt  time.Time            `json:"submittedAt"`
	CompletedAt  *time.Time           `json:"completedAt,omitempty"`
	Error        string               `json:"error,omitempty"`
	ErrorMessage *string              `json:"errorMessage,omitempty"`
	CreatedAt    time.Time            `json:"createdAt"`
}

type BankSubmissionStatus string

const (
	SubmissionStatusSuccess BankSubmissionStatus = "SUCCESS"
	SubmissionStatusFailed  BankSubmissionStatus = "FAILED"
)
