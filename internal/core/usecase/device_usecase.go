package usecase

import (
	"project-home-iot/internal/core/domain"
)

type DeviceUsecase interface {
	RegisterDevice(device *domain.Device) error
	ListDevices() ([]*domain.DeviceSummary, error)
	GetDevice(id string) (*domain.DeviceSummary, error)
	UpdateHeartbeat(id string) error
	UpdateDevice(id string, name string) error
	PairDevice(id string, deviceKey string) error
	UnpairDevice(id string) error
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

func (u *deviceUsecase) GetDevice(id string) (*domain.DeviceSummary, error) {
	return u.repo.GetSummaryByID(id)
}

func (u *deviceUsecase) UpdateDevice(id string, name string) error {
	return u.repo.UpdateName(id, name)
}

func (u *deviceUsecase) PairDevice(id string, deviceKey string) error {
	device, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := u.pairCommander.RequestPair(device.DeviceID, deviceKey); err != nil {
		return err
	}

	caps := DeviceCapabilityMap[device.DeviceType]
	for _, capType := range caps {
		cap, err := u.capabilityRepo.FindByType(capType)
		if err != nil {
			return err
		}
		widget := &domain.Widget{
			DeviceID:      device.DeviceID,
			CapabilityID:  cap.ID,
			WidgetStatus: "exclude",
		}

		if err := u.widgetRepo.CreateWidget(widget); err != nil {
			return err
		}
	}

	return u.pairCommander.Subscribe(device.DeviceID)
}

func (u *deviceUsecase) UnpairDevice(id string) error {
	return u.repo.Unpair(id)
}

func (u *deviceUsecase) UpdateHeartbeat(id string) error {
	return u.repo.UpdateHeartbeat(id)
}