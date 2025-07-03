package services

import (
	"testing"

	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolidBankService_HandleProcessedApplication(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	logEntry := logrus.NewEntry(logger)

	service := &solidBankService{}

	tests := []struct {
		name            string
		application     dto.SolidBankApplication
		bankName        string
		expectedError   string
		expectedStatus  dto.OfferStatus
		shouldHaveOffer bool
	}{
		{
			name: "application with offer should return approved offer",
			application: dto.SolidBankApplication{
				ID:     "test-id-1",
				Status: "PROCESSED",
				Offer: &dto.SolidBankOffer{
					MonthlyPaymentAmount: 500.0,
					TotalRepaymentAmount: 6000.0,
					NumberOfPayments:     12,
					AnnualPercentageRate: 12.5,
					FirstRepaymentDate:   "2024-02-01",
				},
			},
			bankName:        "SolidBank",
			expectedStatus:  dto.OfferStatusApproved,
			shouldHaveOffer: true,
		},
		{
			name: "application without offer should return rejected offer",
			application: dto.SolidBankApplication{
				ID:     "test-id-2",
				Status: "PROCESSED",
				Offer:  nil,
			},
			bankName:        "SolidBank",
			expectedStatus:  dto.OfferStatusRejected,
			shouldHaveOffer: true,
		},
		{
			name: "application with non-PROCESSED status should fail mapping",
			application: dto.SolidBankApplication{
				ID:     "test-id-4",
				Status: "PENDING",
				Offer: &dto.SolidBankOffer{
					MonthlyPaymentAmount: 400.0,
					TotalRepaymentAmount: 4800.0,
					NumberOfPayments:     12,
					AnnualPercentageRate: 11.0,
					FirstRepaymentDate:   "2024-04-01",
				},
			},
			bankName:        "SolidBank",
			expectedError:   "failed to map SolidBank application",
			shouldHaveOffer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offer, err := service.handleProcessedApplication(tt.application, logEntry, tt.bankName)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, offer)
				return
			}

			require.NoError(t, err)

			if tt.shouldHaveOffer {
				require.NotNil(t, offer)
				assert.Equal(t, tt.bankName, offer.BankName)
				assert.Equal(t, tt.expectedStatus, offer.Status)
				assert.NotEmpty(t, offer.ID)
				assert.False(t, offer.CreatedAt.IsZero())

				if tt.expectedStatus == dto.OfferStatusApproved && tt.application.Offer != nil {
					require.NotNil(t, offer.MonthlyPaymentAmount)
					require.NotNil(t, offer.TotalRepaymentAmount)
					require.NotNil(t, offer.NumberOfPayments)
					require.NotNil(t, offer.AnnualPercentageRate)
					require.NotNil(t, offer.FirstRepaymentDate)

					assert.Equal(t, tt.application.Offer.MonthlyPaymentAmount, *offer.MonthlyPaymentAmount)
					assert.Equal(t, tt.application.Offer.TotalRepaymentAmount, *offer.TotalRepaymentAmount)
					assert.Equal(t, tt.application.Offer.NumberOfPayments, *offer.NumberOfPayments)
					assert.Equal(t, tt.application.Offer.AnnualPercentageRate, *offer.AnnualPercentageRate)
					assert.Equal(t, tt.application.Offer.FirstRepaymentDate, *offer.FirstRepaymentDate)
				}

				if tt.expectedStatus == dto.OfferStatusRejected {
					assert.Nil(t, offer.MonthlyPaymentAmount)
					assert.Nil(t, offer.TotalRepaymentAmount)
					assert.Nil(t, offer.NumberOfPayments)
					assert.Nil(t, offer.AnnualPercentageRate)
					assert.Nil(t, offer.FirstRepaymentDate)
				}
			} else {
				assert.Nil(t, offer)
			}
		})
	}
}

func TestSolidBankService_HandleProcessedApplication_EdgeCases(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	logEntry := logrus.NewEntry(logger)
	service := &solidBankService{}

	t.Run("zero values in offer should be preserved", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer: &dto.SolidBankOffer{
				MonthlyPaymentAmount: 0.0,
				TotalRepaymentAmount: 0.0,
				NumberOfPayments:     0,
				AnnualPercentageRate: 0.0,
				FirstRepaymentDate:   "",
			},
		}

		offer, err := service.handleProcessedApplication(app, logEntry, "SolidBank")
		require.NoError(t, err)
		require.NotNil(t, offer)
		assert.Equal(t, dto.OfferStatusApproved, offer.Status)

		require.NotNil(t, offer.MonthlyPaymentAmount)
		require.NotNil(t, offer.TotalRepaymentAmount)
		require.NotNil(t, offer.NumberOfPayments)
		require.NotNil(t, offer.AnnualPercentageRate)
		require.NotNil(t, offer.FirstRepaymentDate)

		assert.Equal(t, 0.0, *offer.MonthlyPaymentAmount)
		assert.Equal(t, 0.0, *offer.TotalRepaymentAmount)
		assert.Equal(t, 0, *offer.NumberOfPayments)
		assert.Equal(t, 0.0, *offer.AnnualPercentageRate)
		assert.Equal(t, "", *offer.FirstRepaymentDate)
	})

	t.Run("application with status other than PROCESSED should fail", func(t *testing.T) {
		testCases := []string{"PENDING", "PROCESSING", "FAILED", "CANCELLED", ""}

		for _, status := range testCases {
			t.Run("status_"+status, func(t *testing.T) {
				app := dto.SolidBankApplication{
					ID:     "test-id",
					Status: status,
					Offer: &dto.SolidBankOffer{
						MonthlyPaymentAmount: 500.0,
						TotalRepaymentAmount: 6000.0,
						NumberOfPayments:     12,
						AnnualPercentageRate: 12.5,
						FirstRepaymentDate:   "2024-02-01",
					},
				}

				offer, err := service.handleProcessedApplication(app, logEntry, "SolidBank")
				require.Error(t, err)
				assert.Contains(t, err.Error(), "failed to map SolidBank application")
				assert.Nil(t, offer)
			})
		}
	})
}

func TestSolidBankService_HandleProcessedApplication_LogMessages(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logEntry := logrus.NewEntry(logger)
	service := &solidBankService{}

	t.Run("approved application should log correct message", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer: &dto.SolidBankOffer{
				MonthlyPaymentAmount: 500.0,
				TotalRepaymentAmount: 6000.0,
				NumberOfPayments:     12,
				AnnualPercentageRate: 12.5,
				FirstRepaymentDate:   "2024-02-01",
			},
		}

		offer, err := service.handleProcessedApplication(app, logEntry, "SolidBank")
		require.NoError(t, err)
		require.NotNil(t, offer)
		assert.Equal(t, dto.OfferStatusApproved, offer.Status)
	})

	t.Run("rejected application should log correct message", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer:  nil,
		}

		offer, err := service.handleProcessedApplication(app, logEntry, "SolidBank")
		require.NoError(t, err)
		require.NotNil(t, offer)
		assert.Equal(t, dto.OfferStatusRejected, offer.Status)
	})
}

func TestSolidBankService_HandleProcessedApplication_OfferValues(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	logEntry := logrus.NewEntry(logger)
	service := &solidBankService{}

	t.Run("high value offer should work", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer: &dto.SolidBankOffer{
				MonthlyPaymentAmount: 9999.99,
				TotalRepaymentAmount: 119999.88,
				NumberOfPayments:     120,
				AnnualPercentageRate: 25.99,
				FirstRepaymentDate:   "2024-12-31",
			},
		}

		offer, err := service.handleProcessedApplication(app, logEntry, "SolidBank")
		require.NoError(t, err)
		require.NotNil(t, offer)
		assert.Equal(t, dto.OfferStatusApproved, offer.Status)

		assert.Equal(t, 9999.99, *offer.MonthlyPaymentAmount)
		assert.Equal(t, 119999.88, *offer.TotalRepaymentAmount)
		assert.Equal(t, 120, *offer.NumberOfPayments)
		assert.Equal(t, 25.99, *offer.AnnualPercentageRate)
		assert.Equal(t, "2024-12-31", *offer.FirstRepaymentDate)
	})

	t.Run("small value offer should work", func(t *testing.T) {
		app := dto.SolidBankApplication{
			ID:     "test-id",
			Status: "PROCESSED",
			Offer: &dto.SolidBankOffer{
				MonthlyPaymentAmount: 0.01,
				TotalRepaymentAmount: 0.12,
				NumberOfPayments:     1,
				AnnualPercentageRate: 0.01,
				FirstRepaymentDate:   "2024-01-01",
			},
		}

		offer, err := service.handleProcessedApplication(app, logEntry, "SolidBank")
		require.NoError(t, err)
		require.NotNil(t, offer)
		assert.Equal(t, dto.OfferStatusApproved, offer.Status)

		assert.Equal(t, 0.01, *offer.MonthlyPaymentAmount)
		assert.Equal(t, 0.12, *offer.TotalRepaymentAmount)
		assert.Equal(t, 1, *offer.NumberOfPayments)
		assert.Equal(t, 0.01, *offer.AnnualPercentageRate)
		assert.Equal(t, "2024-01-01", *offer.FirstRepaymentDate)
	})
}
