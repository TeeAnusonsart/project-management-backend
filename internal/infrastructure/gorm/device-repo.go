package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/mappers"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) CreateDevice(device *domain.Device) error {
	model := mappers.ToDeviceModel(device)
	return r.db.Create(model).Error
}

func (r *DeviceRepository) CreateMonitorData(data *domain.MonitorData) error {
	return r.db.Create(data).Error
}
