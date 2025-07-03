package mappers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBankSubmissionModel(t *testing.T) {
	now := time.Now()
	submissionID := uuid.New()
	testMessage := "Test error message"

	tests := []struct {
		name     string
		input    *dto.BankSubmission
		expected *models.BankSubmission
	}{
		{
			name:     "nil input should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "complete bank submission should map correctly",
			input: &dto.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       dto.SubmissionStatusSuccess,
				BankID:       "bank-123",
				SubmittedAt:  now,
				CompletedAt:  &now,
				Error:        "test-error",
				ErrorMessage: &testMessage,
				CreatedAt:    now,
			},
			expected: &models.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       "SUCCESS",
				BankID:       &[]string{"bank-123"}[0],
				SubmittedAt:  &now,
				CompletedAt:  &now,
				Error:        &[]string{"test-error"}[0],
				ErrorMessage: &testMessage,
				CreatedAt:    now,
			},
		},
		{
			name: "bank submission with empty bankID and error should have nil pointers",
			input: &dto.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       dto.SubmissionStatusSuccess,
				BankID:       "",
				SubmittedAt:  now,
				CompletedAt:  nil,
				Error:        "",
				ErrorMessage: nil,
				CreatedAt:    now,
			},
			expected: &models.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       "SUCCESS",
				BankID:       nil,
				SubmittedAt:  &now,
				CompletedAt:  nil,
				Error:        nil,
				ErrorMessage: nil,
				CreatedAt:    now,
			},
		},
		{
			name: "bank submission with failed status",
			input: &dto.BankSubmission{
				ID:           submissionID,
				BankName:     "FailedBank",
				Status:       dto.SubmissionStatusFailed,
				BankID:       "bank-456",
				SubmittedAt:  now,
				CompletedAt:  &now,
				Error:        "connection-timeout",
				ErrorMessage: &[]string{"Connection timeout occurred"}[0],
				CreatedAt:    now,
			},
			expected: &models.BankSubmission{
				ID:           submissionID,
				BankName:     "FailedBank",
				Status:       "FAILED",
				BankID:       &[]string{"bank-456"}[0],
				SubmittedAt:  &now,
				CompletedAt:  &now,
				Error:        &[]string{"connection-timeout"}[0],
				ErrorMessage: &[]string{"Connection timeout occurred"}[0],
				CreatedAt:    now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBankSubmissionModel(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.BankName, result.BankName)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.BankID, result.BankID)
			assert.Equal(t, tt.expected.SubmittedAt, result.SubmittedAt)
			assert.Equal(t, tt.expected.CompletedAt, result.CompletedAt)
			assert.Equal(t, tt.expected.Error, result.Error)
			assert.Equal(t, tt.expected.ErrorMessage, result.ErrorMessage)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestToBankSubmissionFromModel(t *testing.T) {
	now := time.Now()
	submissionID := uuid.New()
	testMessage := "Test error message"

	tests := []struct {
		name     string
		input    *models.BankSubmission
		expected *dto.BankSubmission
	}{
		{
			name:     "nil input should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "complete bank submission model should map correctly",
			input: &models.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       "SUCCESS",
				BankID:       &[]string{"bank-123"}[0],
				SubmittedAt:  &now,
				CompletedAt:  &now,
				Error:        &[]string{"test-error"}[0],
				ErrorMessage: &testMessage,
				CreatedAt:    now,
			},
			expected: &dto.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       dto.SubmissionStatusSuccess,
				BankID:       "bank-123",
				SubmittedAt:  now,
				CompletedAt:  &now,
				Error:        "test-error",
				ErrorMessage: &testMessage,
				CreatedAt:    now,
			},
		},
		{
			name: "bank submission model with nil pointers should have empty strings",
			input: &models.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       "SUCCESS",
				BankID:       nil,
				SubmittedAt:  &now,
				CompletedAt:  nil,
				Error:        nil,
				ErrorMessage: nil,
				CreatedAt:    now,
			},
			expected: &dto.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       dto.SubmissionStatusSuccess,
				BankID:       "",
				SubmittedAt:  now,
				CompletedAt:  nil,
				Error:        "",
				ErrorMessage: nil,
				CreatedAt:    now,
			},
		},
		{
			name: "bank submission model with nil submitted at",
			input: &models.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       "SUCCESS",
				BankID:       &[]string{"bank-789"}[0],
				SubmittedAt:  nil,
				CompletedAt:  nil,
				Error:        nil,
				ErrorMessage: nil,
				CreatedAt:    now,
			},
			expected: &dto.BankSubmission{
				ID:           submissionID,
				BankName:     "TestBank",
				Status:       dto.SubmissionStatusSuccess,
				BankID:       "bank-789",
				SubmittedAt:  time.Time{},
				CompletedAt:  nil,
				Error:        "",
				ErrorMessage: nil,
				CreatedAt:    now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBankSubmissionFromModel(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.BankName, result.BankName)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.BankID, result.BankID)
			assert.Equal(t, tt.expected.SubmittedAt, result.SubmittedAt)
			assert.Equal(t, tt.expected.CompletedAt, result.CompletedAt)
			assert.Equal(t, tt.expected.Error, result.Error)
			assert.Equal(t, tt.expected.ErrorMessage, result.ErrorMessage)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestBankSubmissionMappers_RoundTrip(t *testing.T) {
	now := time.Now()
	submissionID := uuid.New()
	testMessage := "Test error message"

	t.Run("dto to model to dto should preserve values", func(t *testing.T) {
		original := &dto.BankSubmission{
			ID:           submissionID,
			BankName:     "TestBank",
			Status:       dto.SubmissionStatusSuccess,
			BankID:       "bank-123",
			SubmittedAt:  now,
			CompletedAt:  &now,
			Error:        "test-error",
			ErrorMessage: &testMessage,
			CreatedAt:    now,
		}

		model := ToBankSubmissionModel(original)
		result := ToBankSubmissionFromModel(model)

		require.NotNil(t, result)
		assert.Equal(t, original.ID, result.ID)
		assert.Equal(t, original.BankName, result.BankName)
		assert.Equal(t, original.Status, result.Status)
		assert.Equal(t, original.BankID, result.BankID)
		assert.Equal(t, original.SubmittedAt, result.SubmittedAt)
		assert.Equal(t, original.CompletedAt, result.CompletedAt)
		assert.Equal(t, original.Error, result.Error)
		assert.Equal(t, original.ErrorMessage, result.ErrorMessage)
		assert.Equal(t, original.CreatedAt, result.CreatedAt)
	})

	t.Run("model to dto to model should preserve values", func(t *testing.T) {
		original := &models.BankSubmission{
			ID:           submissionID,
			BankName:     "TestBank",
			Status:       "SUCCESS",
			BankID:       &[]string{"bank-123"}[0],
			SubmittedAt:  &now,
			CompletedAt:  &now,
			Error:        &[]string{"test-error"}[0],
			ErrorMessage: &testMessage,
			CreatedAt:    now,
		}

		dto := ToBankSubmissionFromModel(original)
		result := ToBankSubmissionModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, original.ID, result.ID)
		assert.Equal(t, original.BankName, result.BankName)
		assert.Equal(t, original.Status, result.Status)
		assert.Equal(t, original.BankID, result.BankID)
		assert.Equal(t, original.SubmittedAt, result.SubmittedAt)
		assert.Equal(t, original.CompletedAt, result.CompletedAt)
		assert.Equal(t, original.Error, result.Error)
		assert.Equal(t, original.ErrorMessage, result.ErrorMessage)
		assert.Equal(t, original.CreatedAt, result.CreatedAt)
	})
}
