package mappers

import (
	"gorm.io/gorm"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func WidgetDomainToModel(d *domain.Widget) *models.Widget {
	return &models.Widget{
		Model: gorm.Model{ID: d.ID},
		DeviceID:      d.DeviceID,
		CapabilityID:  d.CapabilityID,
		Value:         d.Value,
		WidgetStatus: d.WidgetStatus,
		WidgetOrder:   d.WidgetOrder,
	}
}

func WidgetModelToDomain(m *models.Widget) *domain.Widget {

	widget := &domain.Widget{
		ID:            m.ID,
		DeviceID:      m.DeviceID,
		CapabilityID:  m.CapabilityID,
		Value:         m.Value,
		WidgetStatus:  m.WidgetStatus,
		WidgetOrder:   m.WidgetOrder,
	}

	if m.Device.DeviceID != "" {
		widget.Device = &domain.Device{
			DeviceID:         m.Device.DeviceID,
			LastHeartbeat: m.Device.LastHeartbeat,
			DeviceName: m.Device.DeviceName,
			DeviceType: m.Device.DeviceType,

		}
	}

	if m.Capability.ID != 0 {
		widget.Capability = &domain.Capability{
			ID:             m.Capability.ID,
			CapabilityType: m.Capability.CapabilityType,
		}
	}

	return widget
}


