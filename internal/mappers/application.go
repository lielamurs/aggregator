package mappers

import (
	"time"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/models"
)

func ToApplicationModel(customerApp *dto.CustomerApplication) *models.Application {
	if customerApp == nil {
		return nil
	}

	app := &models.Application{
		ID:              customerApp.ID,
		Phone:           customerApp.CustomerData.Phone,
		Email:           customerApp.CustomerData.Email,
		MonthlyIncome:   customerApp.CustomerData.MonthlyIncome,
		MonthlyExpenses: customerApp.CustomerData.MonthlyExpenses,
		MaritalStatus:   customerApp.CustomerData.MaritalStatus,
		AgreeToBeScored: customerApp.CustomerData.AgreeToBeScored,
		Amount:          customerApp.CustomerData.Amount,
		Dependents:      customerApp.CustomerData.Dependents,
		Status:          string(customerApp.Status),
		CreatedAt:       customerApp.CreatedAt,
		UpdatedAt:       customerApp.UpdatedAt,
	}

	if len(customerApp.Offers) > 0 {
		app.Offers = make([]models.Offer, len(customerApp.Offers))
		for i, offer := range customerApp.Offers {
			if offerModel := ToOfferModel(&offer); offerModel != nil {
				app.Offers[i] = *offerModel
			}
		}
	}

	if len(customerApp.BankSubmissions) > 0 {
		app.BankSubmissions = make([]models.BankSubmission, len(customerApp.BankSubmissions))
		for i, submission := range customerApp.BankSubmissions {
			if submissionModel := ToBankSubmissionModel(&submission); submissionModel != nil {
				app.BankSubmissions[i] = *submissionModel
			}
		}
	}

	return app
}

func ToCustomerApplicationFromRequest(req *dto.ApplicationRequest) *dto.CustomerApplication {
	if req == nil {
		return nil
	}

	return &dto.CustomerApplication{
		ID:              uuid.New(),
		CustomerData:    *req,
		Status:          dto.StatusPending,
		Offers:          []dto.Offer{},
		BankSubmissions: []dto.BankSubmission{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func ToApplicationStatusResponseFromModel(application *models.Application) *dto.ApplicationStatusResponse {
	if application == nil {
		return nil
	}

	response := &dto.ApplicationStatusResponse{
		ID:        application.ID,
		Status:    dto.ApplicationStatus(application.Status),
		CreatedAt: application.CreatedAt,
		UpdatedAt: application.UpdatedAt,
	}

	if len(application.Offers) > 0 {
		response.Offers = make([]dto.Offer, len(application.Offers))
		for i, offer := range application.Offers {
			if offerDTO := ToOfferFromModel(&offer); offerDTO != nil {
				response.Offers[i] = *offerDTO
			}
		}
	}

	if len(application.BankSubmissions) > 0 {
		response.BankSubmissions = make([]dto.BankSubmission, len(application.BankSubmissions))
		for i, submission := range application.BankSubmissions {
			if submissionDTO := ToBankSubmissionFromModel(&submission); submissionDTO != nil {
				response.BankSubmissions[i] = *submissionDTO
			}
		}
	}

	return response
}
