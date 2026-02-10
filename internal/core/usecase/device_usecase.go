package usecase

import (

	"project-home-iot/internal/core/domain"
)

type DeviceUsecase interface {
	RegisterDevice(device *domain.Device) error
	ListDevices() ([]*domain.DeviceSummary, error)
	GetDevice(id uint) (*domain.DeviceSummary, error)
	UpdateDevice(id uint, name string) error
	PairDevice(id uint, deviceKey string) error
	UnpairDevice(id uint) error
}

type deviceUsecase struct {
	repo          domain.DeviceRepository
	capabilityRepo domain.CapabilityRepository
	widgetRepo     domain.WidgetRepository
	pairCommander  domain.PairCommander
}

func NewDeviceUsecase(
	r domain.DeviceRepository,
	cr domain.CapabilityRepository,
	wr domain.WidgetRepository,
	pc domain.PairCommander,
) DeviceUsecase {
	return &deviceUsecase{
		repo:           r,
		capabilityRepo: cr,
		widgetRepo:     wr,
		pairCommander:  pc,
	}
}

func (u *deviceUsecase) RegisterDevice(device *domain.Device) error {
	// if device.DeviceID == "" {
	// 	return fmt.Errorf("device id required")
	// }
	if err := device.Validate(); err != nil {
		return err
	}

	if err := u.repo.CreateDevice(device); err != nil {
        return err
    }

	return nil
	// return u.repo.CreateDevice(device)
}

// func (u *deviceUsecase) RecordSensorData(data *domain.MonitorData) error {
// 	return u.repo.CreateMonitorData(data)
// }

func (u *deviceUsecase) ListDevices() ([]*domain.DeviceSummary, error) {
	return u.repo.GetAllSummaries()
}

func (u *deviceUsecase) GetDevice(id uint) (*domain.DeviceSummary, error) {
	return u.repo.GetSummaryByID(id)
}

func (u *deviceUsecase) UpdateDevice(id uint, name string) error {
	return u.repo.UpdateName(id, name)
}

func (u *deviceUsecase) PairDevice(id uint, deviceKey string) error {
	device, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := u.pairCommander.RequestPair(device.ID, deviceKey); err != nil {
		return err
	}

	caps := DeviceCapabilityMap[device.DeviceType]
	for _, capType := range caps {
		cap, err := u.capabilityRepo.FindByType(capType)
		if err != nil {
			return err
		}
		widget := &domain.Widget{
			DeviceID:      device.ID,
			CapabilityID:  cap.ID,
			WidgetStatus: "inactive",
		}

		if err := u.widgetRepo.CreateWidget(widget); err != nil {
			return err
		}
	}

	return u.pairCommander.Subscribe(device.Topic)
}

func (u *deviceUsecase) UnpairDevice(id uint) error {
	return u.repo.Unpair(id)
}