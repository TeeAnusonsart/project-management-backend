package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func WidgetDomainToModel(d *domain.Widget) *models.Widget {
	return &models.Widget{
		DeviceID:      d.DeviceID,
		CapabilityID:  d.CapabilityID,
		Value:         d.Value,
		Widget_status: d.Widget_status,
	}
}

func WidgetModelToDomain(m *models.Widget) *domain.Widget {
	return &domain.Widget{
		DeviceID:      m.DeviceID,
		CapabilityID:  m.CapabilityID,
		Value:         m.Value,
		Widget_status: m.Widget_status,
	}
}
