package usecase

import "project-home-iot/internal/core/domain"

type CapabilityUsecase interface {
	// CreateCapability(deviceID string, capabilityType string) error
	FindByType(capType string) (*domain.Capability, error)
}

type capabilityUsecase struct {
	repo domain.CapabilityRepository
}

func NewCapabilityUsecase(r domain.CapabilityRepository) CapabilityUsecase {
	return &capabilityUsecase{repo: r}
}

// func (u *capabilityUsecase) CreateCapability(deviceID string, capabilityType string) error {
// 	device := &domain.Device{
// 		DeviceID: deviceID,
// 	}
// 	return u.repo.CreateCapability(device)
// }

func (u *capabilityUsecase) FindByType(capType string) (*domain.Capability, error) {
	return u.repo.FindByType(capType)
}
