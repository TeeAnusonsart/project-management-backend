package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/mappers"

	"gorm.io/gorm"
)

type WidgetRepository struct {
	db *gorm.DB
}

func NewWidgetRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) CreateWidget(widget *domain.Widget) error {
	model := mappers.WidgetDomainToModel(widget)
	return r.db.Create(model).Error
}
