package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mqtt/dtos"
)

func DomainToDeviceModel(d *domain.Device) *models.Device {
	return &models.Device{
		ID:   d.ID,
		DeviceName: d.DeviceName,
		DeviceType: d.DeviceType,
		Topic:      d.Topic,
	}
}

func ModelToDomainDevice(m *models.Device) *domain.Device {
	return &domain.Device{
		ID:         m.ID,
		// DeviceID:   m.DeviceID,
		DeviceName: m.DeviceName,
		DeviceType: m.DeviceType,
		Topic:      m.Topic,
	}
}

func DomainToDevicePayload(d *domain.Device) *dtos.DevicePayload {
	return &dtos.DevicePayload{
		// DeviceID:   d.DeviceID,
		DeviceName: d.DeviceName,
		DeviceType: d.DeviceType,
		Topic:      d.Topic,
	}
}

func PayloadToDomain(p *dtos.DevicePayload) *domain.Device {
	return &domain.Device{
		ID:         p.DeviceID,
		DeviceName: p.DeviceName,
		DeviceType: p.DeviceType,
		Topic:      p.Topic,
	}
}
