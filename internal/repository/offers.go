package repository

import (
	"context"

	"github.com/lielamurs/aggregator/internal/models"
	"gorm.io/gorm"
)

type OffersRepository struct {
	db *gorm.DB
}

func NewOffersRepository(db *gorm.DB) *OffersRepository {
	return &OffersRepository{
		db: db,
	}
}

func (r *OffersRepository) Create(ctx context.Context, offer *models.Offer) error {
	return r.db.WithContext(ctx).Create(offer).Error
}
