package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func LogDomainToModel(d *domain.Log) *models.Log {
	return &models.Log{
		WidgetID:  d.WidgetID,
		ActorType: d.ActorType,
		Value:     d.Value,
		EventType: d.EventType,
	}
}

func LogModelToDomain(m *models.Log) *domain.Log {
	return &domain.Log{
		ID:        m.ID,
		WidgetID:  m.WidgetID,
		ActorType: m.ActorType,
		Value:     m.Value,
		EventType: m.EventType,
	}
}