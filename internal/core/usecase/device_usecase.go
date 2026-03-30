package usecase

import (
	"project-home-iot/internal/core/domain"
)

type DeviceUsecase interface {
	RegisterDevice(device *domain.Device) error
	ListDevices() ([]*domain.DeviceSummary, error)
	GetUnpairDevice() ([]*domain.DeviceSummary, error)
	GetPairedDevice() ([]*domain.DeviceSummary, error)
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

func (u *deviceUsecase) GetUnpairDevice() ([]*domain.DeviceSummary,error) {
	return u.repo.GetUnpairDevice()
}

func (u *deviceUsecase) GetPairedDevice() ([]*domain.DeviceSummary,error) {
	return u.repo.GetPairedDevice()
}

func (u *deviceUsecase) UpdateDevice(id string, name string) error {
	err := u.repo.UpdateName(id, name)
	if err != nil {
		if err == domain.ErrDeviceNotFound {
			return domain.ErrDeviceNotFound
		}
		return err
	}
	return nil
}

func (u *deviceUsecase) PairDevice(id string, deviceKey string) error {
    device, err := u.repo.GetByID(id)
    if err != nil {
        return err
    }

    exists, err := u.widgetRepo.ExistsByDeviceID(device.DeviceID)
    if err != nil {
        return err
    }
    if exists {
        return domain.ErrDeviceAlreadyPaired
    }
	 

    if err := u.pairCommander.RequestPair(device.DeviceID, deviceKey); err != nil {
        if err.Error() == "invalid device key" {
			 
			return domain.ErrInvalidDeviceKey
		}else if err.Error() == "pair failed" {
			return domain.ErrPairFailed
		}
    }

    capsRefs := DeviceCapabilityMap[device.DeviceType]

    for _, ref := range capsRefs {
        cap, err := u.capabilityRepo.FindByTypeAndControl(ref.Type, ref.Control)
        if err != nil {
            continue
        }

        widget := &domain.Widget{
            DeviceID:     device.DeviceID,
            CapabilityID: cap.ID,
            WidgetStatus: "exclude",
        }

        if err := u.widgetRepo.CreateWidget(widget); err != nil {
            return err
        }
    }

    return u.pairCommander.Subscribe(device.DeviceID)
}

func (u *deviceUsecase) UnpairDevice(id string) error {

	_, err := u.repo.GetByID(id)
    if err != nil {
		if err == domain.ErrDeviceNotFound {
			return domain.ErrDeviceNotFound
		}
        return err
    }

	if err := u.widgetRepo.DeleteByDeviceId(id); err != nil {
		return err
	}
	if err := u.repo.Unpair(id); err != nil {
		return err
	}
	return nil
}

func (u *deviceUsecase) UpdateHeartbeat(id string) error {
	return u.repo.UpdateHeartbeat(id)
}