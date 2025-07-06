package repository

import (
	"context"

	"github.com/lielamurs/aggregator/internal/models"
	"gorm.io/gorm"
)

type BankSubmissionsRepository struct {
	db *gorm.DB
}

func NewBankSubmissionsRepository(db *gorm.DB) *BankSubmissionsRepository {
	return &BankSubmissionsRepository{
		db: db,
	}
}

func (r *BankSubmissionsRepository) Create(ctx context.Context, submission *models.BankSubmission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *BankSubmissionsRepository) Update(ctx context.Context, submission *models.BankSubmission) error {
	return r.db.WithContext(ctx).Save(submission).Error
}
