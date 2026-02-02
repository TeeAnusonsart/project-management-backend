package gorm

import (
	"project-home-iot/internal/core/domain"

	"gorm.io/gorm"
)

type DeviceDataRepository struct {
	db *gorm.DB
}

func NewDeviceDataRepository(db *gorm.DB) *DeviceDataRepository {
	return &DeviceDataRepository{db: db}
}

func (r *DeviceRepository) CreateDeviceData(deviceData *domain.MonitorData) error {
	return r.db.Create(deviceData).Error
}
