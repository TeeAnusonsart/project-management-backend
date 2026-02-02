package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func ToDeviceModel(d *domain.Device) *models.DeviceModel {
	return &models.DeviceModel{
		DeviceID:   d.DeviceID,
		DeviceName: d.DeviceName,
		DeviceKey:  d.DeviceKey,
		DeviceType: d.DeviceType,
		Topic:      d.Topic,
		LastOnline: d.LastOnline,
	}
}

func ToDomainDevice(m *models.DeviceModel) *domain.Device {
	return &domain.Device{
		ID:         m.ID,
		DeviceID:   m.DeviceID,
		DeviceName: m.DeviceName,
		DeviceKey:  m.DeviceKey,
		DeviceType: m.DeviceType,
		Topic:      m.Topic,
		LastOnline: m.LastOnline,
	}
}
