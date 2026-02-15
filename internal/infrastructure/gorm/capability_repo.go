package gorm

import (
	"project-home-iot/internal/core/domain"
	// "project-home-iot/internal/infrastructure/mappers"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mappers"

	"gorm.io/gorm"
)

type CapabilityRepository struct {
	db *gorm.DB
}

func NewCapabilityRepository(db *gorm.DB) *CapabilityRepository {
	return &CapabilityRepository{db: db}
}

func (r *CapabilityRepository) FindByType(capabilityType string) (*domain.Capability, error) {
	var model models.Capability 

	err := r.db.
		Where("capability_type = ?", capabilityType).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return mappers.ModelToDomainCapability(&model), nil
}

func (r *CapabilityRepository) FindByTypeAndControl(capabilityType string, controlType string) (*domain.Capability, error) {
	var model models.Capability
	err := r.db.
		Where("capability_type = ? AND control_type = ?", capabilityType, controlType).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return mappers.ModelToDomainCapability(&model), nil
}
