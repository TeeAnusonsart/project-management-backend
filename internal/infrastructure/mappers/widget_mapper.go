package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func WidgetDomainToModel(d *domain.Widget) *models.WidgetModel {
	return &models.WidgetModel{
		DeviceID:      d.DeviceID,
		CapabilityID:  d.CapabilityID,
		Value:         d.Value,
		Widget_status: d.Widget_status,
	}
}
