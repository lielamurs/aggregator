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

func TestToOfferModel(t *testing.T) {
	now := time.Now()
	offerID := uuid.New()

	tests := []struct {
		name     string
		input    *dto.Offer
		expected *models.Offer
	}{
		{
			name:     "nil input should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "complete offer should map correctly",
			input: &dto.Offer{
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
			expected: &models.Offer{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToOfferModel(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.BankName, result.BankName)
			assert.Equal(t, tt.expected.MonthlyPaymentAmount, result.MonthlyPaymentAmount)
			assert.Equal(t, tt.expected.TotalRepaymentAmount, result.TotalRepaymentAmount)
			assert.Equal(t, tt.expected.NumberOfPayments, result.NumberOfPayments)
			assert.Equal(t, tt.expected.AnnualPercentageRate, result.AnnualPercentageRate)
			assert.Equal(t, tt.expected.FirstRepaymentDate, result.FirstRepaymentDate)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestToOfferFromModel(t *testing.T) {
	now := time.Now()
	offerID := uuid.New()

	tests := []struct {
		name     string
		input    *models.Offer
		expected *dto.Offer
	}{
		{
			name:     "nil input should return nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "complete offer model should map correctly",
			input: &models.Offer{
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
			expected: &dto.Offer{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToOfferFromModel(tt.input)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.BankName, result.BankName)
			assert.Equal(t, tt.expected.MonthlyPaymentAmount, result.MonthlyPaymentAmount)
			assert.Equal(t, tt.expected.TotalRepaymentAmount, result.TotalRepaymentAmount)
			assert.Equal(t, tt.expected.NumberOfPayments, result.NumberOfPayments)
			assert.Equal(t, tt.expected.AnnualPercentageRate, result.AnnualPercentageRate)
			assert.Equal(t, tt.expected.FirstRepaymentDate, result.FirstRepaymentDate)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestOfferMappers_RoundTrip(t *testing.T) {
	now := time.Now()
	offerID := uuid.New()

	t.Run("dto to model to dto should preserve values", func(t *testing.T) {
		original := &dto.Offer{
			ID:                   offerID,
			BankName:             "TestBank",
			MonthlyPaymentAmount: &[]float64{500.0}[0],
			TotalRepaymentAmount: &[]float64{6000.0}[0],
			NumberOfPayments:     &[]int{12}[0],
			AnnualPercentageRate: &[]float64{12.5}[0],
			FirstRepaymentDate:   &[]string{"2024-01-01"}[0],
			Status:               dto.OfferStatusApproved,
			CreatedAt:            now,
		}

		model := ToOfferModel(original)
		result := ToOfferFromModel(model)

		require.NotNil(t, result)
		assert.Equal(t, original.ID, result.ID)
		assert.Equal(t, original.BankName, result.BankName)
		assert.Equal(t, original.MonthlyPaymentAmount, result.MonthlyPaymentAmount)
		assert.Equal(t, original.TotalRepaymentAmount, result.TotalRepaymentAmount)
		assert.Equal(t, original.NumberOfPayments, result.NumberOfPayments)
		assert.Equal(t, original.AnnualPercentageRate, result.AnnualPercentageRate)
		assert.Equal(t, original.FirstRepaymentDate, result.FirstRepaymentDate)
		assert.Equal(t, original.Status, result.Status)
		assert.Equal(t, original.CreatedAt, result.CreatedAt)
	})

	t.Run("model to dto to model should preserve values", func(t *testing.T) {
		original := &models.Offer{
			ID:                   offerID,
			BankName:             "TestBank",
			MonthlyPaymentAmount: &[]float64{500.0}[0],
			TotalRepaymentAmount: &[]float64{6000.0}[0],
			NumberOfPayments:     &[]int{12}[0],
			AnnualPercentageRate: &[]float64{12.5}[0],
			FirstRepaymentDate:   &[]string{"2024-01-01"}[0],
			Status:               "APPROVED",
			CreatedAt:            now,
		}

		dto := ToOfferFromModel(original)
		result := ToOfferModel(dto)

		require.NotNil(t, result)
		assert.Equal(t, original.ID, result.ID)
		assert.Equal(t, original.BankName, result.BankName)
		assert.Equal(t, original.MonthlyPaymentAmount, result.MonthlyPaymentAmount)
		assert.Equal(t, original.TotalRepaymentAmount, result.TotalRepaymentAmount)
		assert.Equal(t, original.NumberOfPayments, result.NumberOfPayments)
		assert.Equal(t, original.AnnualPercentageRate, result.AnnualPercentageRate)
		assert.Equal(t, original.FirstRepaymentDate, result.FirstRepaymentDate)
		assert.Equal(t, original.Status, result.Status)
		assert.Equal(t, original.CreatedAt, result.CreatedAt)
	})
}

func TestOfferMappers_StatusConversion(t *testing.T) {
	now := time.Now()
	offerID := uuid.New()

	t.Run("all valid statuses should convert correctly", func(t *testing.T) {
		statusTests := []struct {
			dtoStatus   dto.OfferStatus
			modelStatus string
		}{
			{dto.OfferStatusApproved, "APPROVED"},
			{dto.OfferStatusRejected, "REJECTED"},
		}

		for _, tt := range statusTests {
			t.Run(string(tt.dtoStatus), func(t *testing.T) {
				dtoOffer := &dto.Offer{
					ID:        offerID,
					BankName:  "TestBank",
					Status:    tt.dtoStatus,
					CreatedAt: now,
				}

				modelOffer := ToOfferModel(dtoOffer)
				require.NotNil(t, modelOffer)
				assert.Equal(t, tt.modelStatus, modelOffer.Status)

				backToDto := ToOfferFromModel(modelOffer)
				require.NotNil(t, backToDto)
				assert.Equal(t, tt.dtoStatus, backToDto.Status)
			})
		}
	})
}
