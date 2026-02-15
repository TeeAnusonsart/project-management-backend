package gorm

import (
	"project-home-iot/internal/core/domain"

	"gorm.io/gorm"
)

type Capability struct {
	gorm.Model
	CapabilityType string
	ControlType	string
	Widgets        []Widget `gorm:"foreignKey:CapabilityID"`
}


type CapabilityRepository struct {
	db *gorm.DB
}

func NewCapabilityRepository(db *gorm.DB) *CapabilityRepository {
	return &CapabilityRepository{db: db}
}

func (r *CapabilityRepository) FindByType(capabilityType string) (*domain.Capability, error) {
	var model Capability 

	err := r.db.
		Where("capability_type = ?", capabilityType).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return ModelToDomainCapability(&model), nil
}

func (r *CapabilityRepository) FindByTypeAndControl(capabilityType string, controlType string) (*domain.Capability, error) {
	var model Capability
	err := r.db.
		Where("capability_type = ? AND control_type = ?", capabilityType, controlType).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return ModelToDomainCapability(&model), nil
}

func DomainToCapabilityModel(d *domain.Capability) *Capability {
	return &Capability{
		CapabilityType: d.CapabilityType,
	}
}

func ModelToDomainCapability(m *Capability) *domain.Capability {
	return &domain.Capability{
		ID : m.ID,
		CapabilityType: m.CapabilityType,
	}
}
