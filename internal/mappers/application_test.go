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

func TestToApplicationModel(t *testing.T) {
	now := time.Now()
	customerAppID := uuid.New()
	offerID := uuid.New()
	submissionID := uuid.New()

	tests := []struct {
		name     string
		input    *dto.CustomerApplication
		expected *models.Application
	}{
		{
			name:     "nil input should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "complete customer application should map correctly",
			input: &dto.CustomerApplication{
				ID: customerAppID,
				CustomerData: dto.ApplicationRequest{
					Phone:           "+1234567890",
					Email:           "test@example.com",
					MonthlyIncome:   5000.0,
					MonthlyExpenses: 2000.0,
					MaritalStatus:   "single",
					AgreeToBeScored: true,
					Amount:          10000.0,
					Dependents:      2,
				},
				Status: dto.StatusPending,
				Offers: []dto.Offer{
					{
						ID:                   offerID,
						BankName:             "TestBank",
						MonthlyPaymentAmount: &[]float64{500.0}[0],
						TotalRepaymentAmount: &[]float64{6000.0}[0],
						NumberOfPayments:     &[]int{12}[0],
						AnnualPercentageRate: &[]float64{12.5}[0],
						FirstRepaymentDate:   &[]string{"2024-01-01"}[0],
						Status:               dto.OfferStatusApproved,
						CreatedAt:            now,
					},
				},
				BankSubmissions: []dto.BankSubmission{
					{
						ID:           submissionID,
						BankName:     "TestBank",
						Status:       dto.SubmissionStatusSuccess,
						BankID:       "bank-123",
						SubmittedAt:  now,
						CompletedAt:  &now,
						Error:        "",
						ErrorMessage: nil,
						CreatedAt:    now,
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			expected: &models.Application{
				ID:              customerAppID,
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "single",
				AgreeToBeScored: true,
				Amount:          10000.0,
				Dependents:      2,
				Status:          "PENDING",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
		{
			name: "customer application without offers and submissions",
			input: &dto.CustomerApplication{
				ID: customerAppID,
				CustomerData: dto.ApplicationRequest{
					Phone:           "+1234567890",
					Email:           "test@example.com",
					MonthlyIncome:   3000.0,
					MonthlyExpenses: 1500.0,
					MaritalStatus:   "married",
					AgreeToBeScored: false,
					Amount:          5000.0,
					Dependents:      1,
				},
				Status:          dto.StatusCompleted,
				Offers:          []dto.Offer{},
				BankSubmissions: []dto.BankSubmission{},
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expected: &models.Application{
				ID:              customerAppID,
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   3000.0,
				MonthlyExpenses: 1500.0,
				MaritalStatus:   "married",
				AgreeToBeScored: false,
				Amount:          5000.0,
				Dependents:      1,
				Status:          "COMPLETED",
				CreatedAt:       now,
				UpdatedAt:       now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToApplicationModel(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Phone, result.Phone)
			assert.Equal(t, tt.expected.Email, result.Email)
			assert.Equal(t, tt.expected.MonthlyIncome, result.MonthlyIncome)
			assert.Equal(t, tt.expected.MonthlyExpenses, result.MonthlyExpenses)
			assert.Equal(t, tt.expected.MaritalStatus, result.MaritalStatus)
			assert.Equal(t, tt.expected.AgreeToBeScored, result.AgreeToBeScored)
			assert.Equal(t, tt.expected.Amount, result.Amount)
			assert.Equal(t, tt.expected.Dependents, result.Dependents)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}

func TestToCustomerApplicationFromRequest(t *testing.T) {
	tests := []struct {
		name  string
		input *dto.ApplicationRequest
	}{
		{
			name:  "nil input should return nil",
			input: nil,
		},
		{
			name: "complete application request should map correctly",
			input: &dto.ApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "single",
				AgreeToBeScored: true,
				Amount:          10000.0,
				Dependents:      2,
			},
		},
		{
			name: "minimal application request should map correctly",
			input: &dto.ApplicationRequest{
				Phone:           "+0000000000",
				Email:           "min@example.com",
				MonthlyIncome:   1000.0,
				MonthlyExpenses: 500.0,
				MaritalStatus:   "married",
				AgreeToBeScored: false,
				Amount:          1000.0,
				Dependents:      0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCustomerApplicationFromRequest(tt.input)

			if tt.input == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.NotEqual(t, uuid.Nil, result.ID)
			assert.Equal(t, *tt.input, result.CustomerData)
			assert.Equal(t, dto.StatusPending, result.Status)
			assert.Empty(t, result.Offers)
			assert.Empty(t, result.BankSubmissions)
			assert.False(t, result.CreatedAt.IsZero())
			assert.False(t, result.UpdatedAt.IsZero())
		})
	}
}

func TestToApplicationStatusResponseFromModel(t *testing.T) {
	now := time.Now()
	appID := uuid.New()
	offerID := uuid.New()
	submissionID := uuid.New()

	tests := []struct {
		name     string
		input    *models.Application
		expected *dto.ApplicationStatusResponse
	}{
		{
			name:     "nil input should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "complete application model should map correctly",
			input: &models.Application{
				ID:              appID,
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "single",
				AgreeToBeScored: true,
				Amount:          10000.0,
				Dependents:      2,
				Status:          "PENDING",
				Offers: []models.Offer{
					{
						ID:                   offerID,
						BankName:             "TestBank",
						MonthlyPaymentAmount: &[]float64{500.0}[0],
						TotalRepaymentAmount: &[]float64{6000.0}[0],
						NumberOfPayments:     &[]int{12}[0],
						AnnualPercentageRate: &[]float64{12.5}[0],
						FirstRepaymentDate:   &[]string{"2024-01-01"}[0],
						Status:               "APPROVED",
						CreatedAt:            now,
					},
				},
				BankSubmissions: []models.BankSubmission{
					{
						ID:           submissionID,
						BankName:     "TestBank",
						Status:       "COMPLETED",
						BankID:       &[]string{"bank-123"}[0],
						SubmittedAt:  &now,
						CompletedAt:  &now,
						Error:        nil,
						ErrorMessage: nil,
						CreatedAt:    now,
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
			expected: &dto.ApplicationStatusResponse{
				ID:        appID,
				Status:    dto.StatusPending,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		{
			name: "application model without offers and submissions",
			input: &models.Application{
				ID:              appID,
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   3000.0,
				MonthlyExpenses: 1500.0,
				MaritalStatus:   "married",
				AgreeToBeScored: false,
				Amount:          5000.0,
				Dependents:      1,
				Status:          "COMPLETED",
				Offers:          []models.Offer{},
				BankSubmissions: []models.BankSubmission{},
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			expected: &dto.ApplicationStatusResponse{
				ID:        appID,
				Status:    dto.StatusCompleted,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToApplicationStatusResponseFromModel(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}

func TestToCustomerApplicationFromRequest_EdgeCases(t *testing.T) {
	t.Run("each call should generate unique ID", func(t *testing.T) {
		input := &dto.ApplicationRequest{
			Phone:           "+1234567890",
			Email:           "test@example.com",
			MonthlyIncome:   5000.0,
			MonthlyExpenses: 2000.0,
			MaritalStatus:   "single",
			AgreeToBeScored: true,
			Amount:          10000.0,
			Dependents:      2,
		}

		result1 := ToCustomerApplicationFromRequest(input)
		result2 := ToCustomerApplicationFromRequest(input)

		require.NotNil(t, result1)
		require.NotNil(t, result2)
		assert.NotEqual(t, result1.ID, result2.ID)
	})
}
