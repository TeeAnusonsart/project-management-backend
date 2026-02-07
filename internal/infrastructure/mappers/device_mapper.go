package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mqtt/dtos"
)

func DomainToDeviceModel(d *domain.Device) *models.Device {
	return &models.Device{
		// DeviceID:   d.DeviceID,
		DeviceName: d.DeviceName,
		DeviceKey:  d.DeviceKey,
		DeviceType: d.DeviceType,
		Topic:      d.Topic,
	}
}

func ModelToDomainDevice(m *models.Device) *domain.Device {
	return &domain.Device{
		ID:         m.ID,
		// DeviceID:   m.DeviceID,
		DeviceName: m.DeviceName,
		DeviceKey:  m.DeviceKey,
		DeviceType: m.DeviceType,
		Topic:      m.Topic,
	}
}

func DomainToDevicePayload(d *domain.Device) *dtos.DevicePayload {
	return &dtos.DevicePayload{
		// DeviceID:   d.DeviceID,
		DeviceName: d.DeviceName,
		DeviceType: d.DeviceType,
		DeviceKey:  d.DeviceKey,
		Topic:      d.Topic,
	}
}

func PayloadToDomain(p *dtos.DevicePayload) *domain.Device {
	return &domain.Device{
		// DeviceID:   p.DeviceID,
		DeviceName: p.DeviceName,
		DeviceType: p.DeviceType,
		DeviceKey:  p.DeviceKey,
		Topic:      p.Topic,
	}
}
