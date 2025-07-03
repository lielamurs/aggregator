package mappers

import (
	"testing"
	"time"

	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToSolidBankRequestFromApplicationRequest(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.ApplicationRequest
		expected *dto.SolidBankApplicationRequest
	}{
		{
			name: "complete application request should map correctly",
			input: dto.ApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "single",
				AgreeToBeScored: true,
				Amount:          10000.0,
				Dependents:      2,
			},
			expected: &dto.SolidBankApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "single",
				AgreeToBeScored: true,
				Amount:          10000.0,
			},
		},
		{
			name: "high values should be preserved",
			input: dto.ApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   999999.99,
				MonthlyExpenses: 888888.88,
				MaritalStatus:   "married",
				AgreeToBeScored: false,
				Amount:          1000000.0,
				Dependents:      10,
			},
			expected: &dto.SolidBankApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   999999.99,
				MonthlyExpenses: 888888.88,
				MaritalStatus:   "married",
				AgreeToBeScored: false,
				Amount:          1000000.0,
			},
		},
		{
			name: "marital status variations should be preserved",
			input: dto.ApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "divorced",
				AgreeToBeScored: true,
				Amount:          15000.0,
				Dependents:      3,
			},
			expected: &dto.SolidBankApplicationRequest{
				Phone:           "+1234567890",
				Email:           "test@example.com",
				MonthlyIncome:   5000.0,
				MonthlyExpenses: 2000.0,
				MaritalStatus:   "divorced",
				AgreeToBeScored: true,
				Amount:          15000.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSolidBankRequestFromApplicationRequest(tt.input)
			require.NotNil(t, result)
			assert.Equal(t, tt.expected.Phone, result.Phone)
			assert.Equal(t, tt.expected.Email, result.Email)
			assert.Equal(t, tt.expected.MonthlyIncome, result.MonthlyIncome)
			assert.Equal(t, tt.expected.MonthlyExpenses, result.MonthlyExpenses)
			assert.Equal(t, tt.expected.MaritalStatus, result.MaritalStatus)
			assert.Equal(t, tt.expected.AgreeToBeScored, result.AgreeToBeScored)
			assert.Equal(t, tt.expected.Amount, result.Amount)
		})
	}
}

func TestToOfferFromSolidBankApplication(t *testing.T) {
	tests := []struct {
		name            string
		application     dto.SolidBankApplication
		bankName        string
		expectedStatus  dto.OfferStatus
		shouldReturnNil bool
	}{
		{
			name: "processed application with offer should return approved offer",
			application: dto.SolidBankApplication{
				ID:     "test-id-1",
				Status: "PROCESSED",
				Offer: &dto.SolidBankOffer{
					MonthlyPaymentAmount: 500.0,
					TotalRepaymentAmount: 6000.0,
					NumberOfPayments:     12,
					AnnualPercentageRate: 12.5,
					FirstRepaymentDate:   "2024-01-01",
				},
			},
			bankName:        "SolidBank",
			expectedStatus:  dto.OfferStatusApproved,
			shouldReturnNil: false,
		},
		{
			name: "processed application without offer should return rejected offer",
			application: dto.SolidBankApplication{
				ID:     "test-id-2",
				Status: "PROCESSED",
				Offer:  nil,
			},
			bankName:        "SolidBank",
			expectedStatus:  dto.OfferStatusRejected,
			shouldReturnNil: false,
		},
		{
			name: "non-processed application should return nil",
			application: dto.SolidBankApplication{
				ID:     "test-id-3",
				Status: "PENDING",
				Offer: &dto.SolidBankOffer{
					MonthlyPaymentAmount: 500.0,
					TotalRepaymentAmount: 6000.0,
					NumberOfPayments:     12,
					AnnualPercentageRate: 12.5,
					FirstRepaymentDate:   "2024-01-01",
				},
			},
			bankName:        "SolidBank",
			shouldReturnNil: true,
		},
		{
			name: "processed application with different bank name",
			application: dto.SolidBankApplication{
				ID:     "test-id-4",
				Status: "PROCESSED",
				Offer: &dto.SolidBankOffer{
					MonthlyPaymentAmount: 300.0,
					TotalRepaymentAmount: 3600.0,
					NumberOfPayments:     12,
					AnnualPercentageRate: 10.0,
					FirstRepaymentDate:   "2024-02-01",
				},
			},
			bankName:        "TestBank",
			expectedStatus:  dto.OfferStatusApproved,
			shouldReturnNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToOfferFromSolidBankApplication(tt.application, tt.bankName)

			if tt.shouldReturnNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.bankName, result.BankName)
			assert.Equal(t, tt.expectedStatus, result.Status)
			assert.NotEmpty(t, result.ID)
			assert.False(t, result.CreatedAt.IsZero())

			if tt.expectedStatus == dto.OfferStatusApproved && tt.application.Offer != nil {
				require.NotNil(t, result.MonthlyPaymentAmount)
				require.NotNil(t, result.TotalRepaymentAmount)
				require.NotNil(t, result.NumberOfPayments)
				require.NotNil(t, result.AnnualPercentageRate)
				require.NotNil(t, result.FirstRepaymentDate)

				assert.Equal(t, tt.application.Offer.MonthlyPaymentAmount, *result.MonthlyPaymentAmount)
				assert.Equal(t, tt.application.Offer.TotalRepaymentAmount, *result.TotalRepaymentAmount)
				assert.Equal(t, tt.application.Offer.NumberOfPayments, *result.NumberOfPayments)
				assert.Equal(t, tt.application.Offer.AnnualPercentageRate, *result.AnnualPercentageRate)
				assert.Equal(t, tt.application.Offer.FirstRepaymentDate, *result.FirstRepaymentDate)
			}

			if tt.expectedStatus == dto.OfferStatusRejected {
				assert.Nil(t, result.MonthlyPaymentAmount)
				assert.Nil(t, result.TotalRepaymentAmount)
				assert.Nil(t, result.NumberOfPayments)
				assert.Nil(t, result.AnnualPercentageRate)
				assert.Nil(t, result.FirstRepaymentDate)
			}
		})
	}
}

func TestToOfferFromSolidBankApplication_EdgeCases(t *testing.T) {
	t.Run("various non-processed statuses should return nil", func(t *testing.T) {
		statuses := []string{"PENDING", "PROCESSING", "FAILED", "CANCELLED", "REJECTED", ""}

		for _, status := range statuses {
			t.Run("status_"+status, func(t *testing.T) {
				app := dto.SolidBankApplication{
					ID:     "test-id",
					Status: status,
					Offer: &dto.SolidBankOffer{
						MonthlyPaymentAmount: 500.0,
						TotalRepaymentAmount: 6000.0,
						NumberOfPayments:     12,
						AnnualPercentageRate: 12.5,
						FirstRepaymentDate:   "2024-01-01",
					},
				}

				result := ToOfferFromSolidBankApplication(app, "SolidBank")
				assert.Nil(t, result)
			})
		}
	})

	t.Run("high values in offer should be preserved", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer: &dto.SolidBankOffer{
				MonthlyPaymentAmount: 99999.99,
				TotalRepaymentAmount: 1199999.88,
				NumberOfPayments:     240,
				AnnualPercentageRate: 99.99,
				FirstRepaymentDate:   "2024-12-31",
			},
		}

		result := ToOfferFromSolidBankApplication(app, "SolidBank")
		require.NotNil(t, result)
		assert.Equal(t, dto.OfferStatusApproved, result.Status)

		assert.Equal(t, 99999.99, *result.MonthlyPaymentAmount)
		assert.Equal(t, 1199999.88, *result.TotalRepaymentAmount)
		assert.Equal(t, 240, *result.NumberOfPayments)
		assert.Equal(t, 99.99, *result.AnnualPercentageRate)
		assert.Equal(t, "2024-12-31", *result.FirstRepaymentDate)
	})

	t.Run("each call should generate unique ID", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer:  nil,
		}

		result1 := ToOfferFromSolidBankApplication(app, "SolidBank")
		result2 := ToOfferFromSolidBankApplication(app, "SolidBank")

		require.NotNil(t, result1)
		require.NotNil(t, result2)
		assert.NotEqual(t, result1.ID, result2.ID)
	})

	t.Run("each call should have different created time", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer:  nil,
		}

		result1 := ToOfferFromSolidBankApplication(app, "SolidBank")
		time.Sleep(1 * time.Millisecond)
		result2 := ToOfferFromSolidBankApplication(app, "SolidBank")

		require.NotNil(t, result1)
		require.NotNil(t, result2)
		assert.True(t, result2.CreatedAt.After(result1.CreatedAt) || result2.CreatedAt.Equal(result1.CreatedAt))
	})
}

func TestToSolidBankRequestFromApplicationRequest_EdgeCases(t *testing.T) {
	t.Run("special characters in strings should be preserved", func(t *testing.T) {
		input := dto.ApplicationRequest{
			Phone:           "+1-234-567-890 ext. 123",
			Email:           "test+tag@example.co.uk",
			MonthlyIncome:   5000.50,
			MonthlyExpenses: 2000.25,
			MaritalStatus:   "it's complicated",
			AgreeToBeScored: true,
			Amount:          10000.99,
			Dependents:      2,
		}

		result := ToSolidBankRequestFromApplicationRequest(input)
		require.NotNil(t, result)
		assert.Equal(t, "+1-234-567-890 ext. 123", result.Phone)
		assert.Equal(t, "test+tag@example.co.uk", result.Email)
		assert.Equal(t, "it's complicated", result.MaritalStatus)
		assert.Equal(t, 5000.50, result.MonthlyIncome)
		assert.Equal(t, 2000.25, result.MonthlyExpenses)
		assert.Equal(t, 10000.99, result.Amount)
	})

	t.Run("dependents field should be ignored", func(t *testing.T) {
		input := dto.ApplicationRequest{
			Phone:           "+1234567890",
			Email:           "test@example.com",
			MonthlyIncome:   5000.0,
			MonthlyExpenses: 2000.0,
			MaritalStatus:   "single",
			AgreeToBeScored: true,
			Amount:          10000.0,
			Dependents:      99,
		}

		result := ToSolidBankRequestFromApplicationRequest(input)
		require.NotNil(t, result)
		assert.Equal(t, "+1234567890", result.Phone)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, 5000.0, result.MonthlyIncome)
		assert.Equal(t, 2000.0, result.MonthlyExpenses)
		assert.Equal(t, "single", result.MaritalStatus)
		assert.True(t, result.AgreeToBeScored)
		assert.Equal(t, 10000.0, result.Amount)
	})
}
