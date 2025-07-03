package mappers

import (
	"time"

	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/models"
)

func ToBankSubmissionModel(bankSubmission *dto.BankSubmission) *models.BankSubmission {
	if bankSubmission == nil {
		return nil
	}

	var bankID *string
	if bankSubmission.BankID != "" {
		bankID = &bankSubmission.BankID
	}

	var error *string
	if bankSubmission.Error != "" {
		error = &bankSubmission.Error
	}

	return &models.BankSubmission{
		ID:           bankSubmission.ID,
		BankName:     bankSubmission.BankName,
		Status:       string(bankSubmission.Status),
		BankID:       bankID,
		SubmittedAt:  &bankSubmission.SubmittedAt,
		CompletedAt:  bankSubmission.CompletedAt,
		Error:        error,
		ErrorMessage: bankSubmission.ErrorMessage,
		CreatedAt:    bankSubmission.CreatedAt,
	}
}

func ToBankSubmissionFromModel(bankSubmission *models.BankSubmission) *dto.BankSubmission {
	if bankSubmission == nil {
		return nil
	}

	var bankID string
	if bankSubmission.BankID != nil {
		bankID = *bankSubmission.BankID
	}

	var error string
	if bankSubmission.Error != nil {
		error = *bankSubmission.Error
	}

	var submittedAt time.Time
	if bankSubmission.SubmittedAt != nil {
		submittedAt = *bankSubmission.SubmittedAt
	}

	return &dto.BankSubmission{
		ID:           bankSubmission.ID,
		BankName:     bankSubmission.BankName,
		Status:       dto.BankSubmissionStatus(bankSubmission.Status),
		BankID:       bankID,
		SubmittedAt:  submittedAt,
		CompletedAt:  bankSubmission.CompletedAt,
		Error:        error,
		ErrorMessage: bankSubmission.ErrorMessage,
		CreatedAt:    bankSubmission.CreatedAt,
	}
}
