package usecase

import (
	// "fmt"
	"fmt"
	"project-home-iot/internal/core/domain"
)

type DeviceUsecase interface {
	RegisterDevice(device *domain.Device) error
	// RecordSensorData(data *domain.MonitorData) error
	ListDevices() ([]*domain.Device, error)
	GetDevice(id uint) (*domain.Device, error)
	UpdateDevice(id uint, name string) error
	PairDevice(id uint, deviceKey string) error
	UnpairDevice(id uint) error
}

type deviceUsecase struct {
	repo              domain.DeviceRepository
	capabilityRepo domain.CapabilityRepository
	widgetRepo     domain.WidgetRepository

}

func NewDeviceUsecase(r domain.DeviceRepository, cuc CapabilityUsecase, wuc WidgetUsecase) DeviceUsecase {
	return &deviceUsecase{repo: r, capabilityRepo: cuc, widgetRepo: wuc}
}

func (u *deviceUsecase) RegisterDevice(device *domain.Device) error {
	// if device.DeviceID == "" {
	// 	return fmt.Errorf("device id required")
	// }
	if err := u.repo.CreateDevice(device); err != nil {
        return err
    }

	fmt.Printf("Device created with ID: %d\n", device.ID)
	caps := DeviceCapabilityMap[device.DeviceType]
	for _, capType := range caps {
		cap, err := u.capabilityRepo.FindByType(capType)
		if err != nil {
			return err
		}
		widget := &domain.Widget{
			DeviceID:      device.ID,
			CapabilityID:  cap.ID,
			Widget_status: "inactive",
		}

		if err := u.widgetRepo.CreateWidget(widget); err != nil {
			return err
		}
	}
	return nil
	// return u.repo.CreateDevice(device)
}

// func (u *deviceUsecase) RecordSensorData(data *domain.MonitorData) error {
// 	return u.repo.CreateMonitorData(data)
// }

func (u *deviceUsecase) ListDevices() ([]*domain.Device, error) {
	return u.deviceRepo.GetAll()
}

func (u *deviceUsecase) GetDevice(id uint) (*domain.Device, error) {
	return u.deviceRepo.GetByID(id)
}

func (u *deviceUsecase) UpdateDevice(id uint, name string) error {
	return u.deviceRepo.UpdateName(id, name)
}

func (u *deviceUsecase) PairDevice(id uint, deviceKey string) error {
	return u.deviceRepo.Pair(id, deviceKey)
}

func (u *deviceUsecase) UnpairDevice(id uint) error {
	return u.deviceRepo.Unpair(id)
}