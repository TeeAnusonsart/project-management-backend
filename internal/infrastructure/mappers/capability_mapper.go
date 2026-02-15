package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func DomainToCapabilityModel(d *domain.Capability) *models.Capability {
	return &models.Capability{
		CapabilityType: d.CapabilityType,
	}
}

func ModelToDomainCapability(m *models.Capability) *domain.Capability {
	return &domain.Capability{
		ID : m.ID,
		CapabilityType: m.CapabilityType,
	}
}