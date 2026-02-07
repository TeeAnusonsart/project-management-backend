package usecase

import (
	"fmt"
	"project-home-iot/internal/core/domain"
)

type DeviceUsecase interface {
	RegisterDevice(device *domain.Device) error
	// RecordSensorData(data *domain.MonitorData) error
}

type deviceUsecase struct {
	repo              domain.DeviceRepository
	capabilityUsecase CapabilityUsecase
	widgetUsecase     WidgetUsecase
}

func NewDeviceUsecase(r domain.DeviceRepository, cuc CapabilityUsecase, wuc WidgetUsecase) DeviceUsecase {
	return &deviceUsecase{repo: r, capabilityUsecase: cuc, widgetUsecase: wuc}
}

func (u *deviceUsecase) RegisterDevice(device *domain.Device) error {
	if device.DeviceID == "" {
		return fmt.Errorf("device id required")
	}

	caps := DeviceCapabilityMap[device.DeviceType]

	for _, capType := range caps {

		cap, err := u.capabilityUsecase.FindByType(capType)
		if err != nil {
			return err
		}

		widget := &domain.Widget{
			DeviceID:      device.ID,
			CapabilityID:  cap.ID,
			Widget_status: "inactive",
			Value:         0,
		}

		if err := u.widgetUsecase.CreateWidget(widget); err != nil {
			return err
		}
	}

	return u.repo.CreateDevice(device)
}

// func (u *deviceUsecase) RecordSensorData(data *domain.MonitorData) error {
// 	return u.repo.CreateMonitorData(data)
// }
