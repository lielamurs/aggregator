package models

import (
	"time"

	"github.com/google/uuid"
)

type BankSubmission struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	BankName      string
	Status        string
	BankID        *string
	SubmittedAt   *time.Time
	CompletedAt   *time.Time
	Error         *string
	ErrorMessage  *string
	CreatedAt     time.Time
}
