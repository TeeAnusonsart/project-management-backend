package usecase

import (
	"fmt"
	"project-home-iot/internal/core/domain"
)

type DeviceUsecase interface {
	RegisterDevice(device *domain.Device) error
	RecordSensorData(data *domain.MonitorData) error
}

type deviceUsecase struct {
	repo domain.DeviceRepository
}

func NewDeviceUsecase(r domain.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{repo: r}
}

func (u *deviceUsecase) RegisterDevice(device *domain.Device) error {
	if device.DeviceID == "" {
		return fmt.Errorf("device id required")
	}
	return u.repo.CreateDevice(device)
}

func (u *deviceUsecase) RecordSensorData(data *domain.MonitorData) error {
	return u.repo.CreateMonitorData(data)
}
