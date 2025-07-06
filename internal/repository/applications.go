package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lielamurs/aggregator/internal/models"
	"gorm.io/gorm"
)

type ApplicationsRepository struct {
	db *gorm.DB
}

func NewApplicationsRepository(db *gorm.DB) *ApplicationsRepository {
	return &ApplicationsRepository{
		db: db,
	}
}

func (r *ApplicationsRepository) Create(ctx context.Context, app *models.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *ApplicationsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	var app models.Application
	err := r.db.WithContext(ctx).Preload("Offers").Preload("BankSubmissions").First(&app, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *ApplicationsRepository) Update(ctx context.Context, app *models.Application) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *ApplicationsRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Application{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ApplicationsRepository) GetProcessingApplicationsWithBankSubmissions(ctx context.Context) ([]models.Application, error) {
	var apps []models.Application
	err := r.db.WithContext(ctx).Preload("BankSubmissions").Where("status = ?", "PROCESSING").Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}
