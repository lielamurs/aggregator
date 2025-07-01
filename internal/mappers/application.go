package mappers

import (
	"github.com/lielamurs/aggregator/internal/dto"
)

func ApplicationStatusFromBankStatuses(submissions []dto.BankSubmission) dto.ApplicationStatus {
	if len(submissions) == 0 {
		return dto.StatusPending
	}

	allCompleted := true
	hasError := false

	for _, submission := range submissions {
		switch submission.Status {
		case dto.BankStatusPending, dto.BankStatusSubmitted:
			allCompleted = false
		case dto.BankStatusError:
			hasError = true
		}
	}

	if hasError && allCompleted {
		return dto.StatusFailed
	}

	if allCompleted {
		return dto.StatusCompleted
	}

	hasSubmitted := false
	for _, submission := range submissions {
		if submission.Status != dto.BankStatusPending {
			hasSubmitted = true
			break
		}
	}

	if hasSubmitted {
		return dto.StatusProcessing
	}

	return dto.StatusPending
}
