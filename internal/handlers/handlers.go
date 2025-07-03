package handlers

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lielamurs/aggregator/internal/dto"
	"github.com/lielamurs/aggregator/internal/mappers"
	"github.com/lielamurs/aggregator/internal/services"
	"github.com/sirupsen/logrus"
)

type ApplicationHandler struct {
	applicationService services.ApplicationService
	validator          *validator.Validate
	logger             *logrus.Logger
}

func NewApplicationHandler(applicationService services.ApplicationService, logger *logrus.Logger) *ApplicationHandler {
	return &ApplicationHandler{
		applicationService: applicationService,
		validator:          validator.New(),
		logger:             logger,
	}
}

func (h *ApplicationHandler) SubmitApplication(c echo.Context) error {
	var req dto.ApplicationRequest

	if err := c.Bind(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind application request")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request format",
			Code:    "INVALID_REQUEST_FORMAT",
		})
	}

	if err := h.validator.Struct(&req); err != nil {
		h.logger.WithError(err).WithField("request", req).Error("Application request validation failed")

		validationErrors := extractValidationErrors(err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation Failed",
			Message: validationErrors,
			Code:    "VALIDATION_FAILED",
		})
	}

	app := mappers.ToCustomerApplicationFromRequest(&req)
	response, err := h.applicationService.SubmitApplication(c.Request().Context(), app)
	if err != nil {
		h.logger.WithError(err).Error("Failed to submit application")
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to process application",
			Code:    "APPLICATION_PROCESSING_FAILED",
		})
	}

	h.logger.WithField("application_id", response.ID).Info("Application submitted successfully")

	return c.JSON(http.StatusCreated, response)
}

func (h *ApplicationHandler) GetApplicationStatus(c echo.Context) error {
	id := c.Param("id")
	applicationID, err := uuid.Parse(id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Invalid application ID format")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid application ID format",
			Code:    "INVALID_APPLICATION_ID",
		})
	}

	h.logger.WithField("application_id", applicationID).Info("Retrieving application status")

	modelApp, err := h.applicationService.GetApplicationStatus(c.Request().Context(), applicationID)
	if err != nil {
		h.logger.WithError(err).WithField("application_id", applicationID).Error("Failed to get application status")

		if isNotFoundError(err) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "Application not found",
				Code:    "APPLICATION_NOT_FOUND",
			})
		}

		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve application status",
			Code:    "APPLICATION_RETRIEVAL_FAILED",
		})
	}

	h.logger.WithFields(logrus.Fields{
		"application_id": applicationID,
		"status":         modelApp.Status,
		"offers_count":   len(modelApp.Offers),
	}).Info("Application status retrieved successfully")

	response := mappers.ToApplicationStatusResponseFromModel(modelApp)
	return c.JSON(http.StatusOK, response)
}

func (h *ApplicationHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "healthy",
		"service": "financing-application-aggregator",
	})
}

func extractValidationErrors(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "Validation failed"
	}

	var messages []string
	for _, validationError := range validationErrors {
		switch validationError.Tag() {
		case "required":
			messages = append(messages, validationError.Field()+" is required")
		case "email":
			messages = append(messages, validationError.Field()+" must be a valid email address")
		case "min":
			messages = append(messages, validationError.Field()+" must be greater than or equal to "+validationError.Param())
		case "oneof":
			messages = append(messages, validationError.Field()+" must be one of: "+validationError.Param())
		default:
			messages = append(messages, validationError.Field()+" is invalid")
		}
	}

	if len(messages) == 0 {
		return "Validation failed"
	}

	result := messages[0]
	for i := 1; i < len(messages); i++ {
		result += ", " + messages[i]
	}

	return result
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "does not exist") ||
		strings.Contains(errMsg, "application with id") && strings.Contains(errMsg, "not found")
}
