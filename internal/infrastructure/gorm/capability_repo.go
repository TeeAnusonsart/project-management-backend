package gorm

import (
	"project-home-iot/internal/core/domain"
	// "project-home-iot/internal/infrastructure/mappers"

	"gorm.io/gorm"
)

type CapabilityRepository struct {
	db *gorm.DB
}

func NewCapabilityRepository(db *gorm.DB) *CapabilityRepository {
	return &CapabilityRepository{db: db}
}

func (r *CapabilityRepository) FindByType(capabilityType string) (*domain.Capability, error) {
	capability := domain.Capability{}
	err := r.db.Where("capability_type = ?", capabilityType).First(&capability).Error

	if err != nil {
		return nil, err
	}

	return &capability, nil
}
