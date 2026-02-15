package usecase

import (
	"project-home-iot/internal/core/domain"
)

type CommandUsecase interface {
	SendCommand(cmd *domain.DeviceCommand, widgetID uint, correlationID string) error
}

type commandUsecase struct {
	deviceRepo       domain.DeviceRepository
	widgetRepository domain.WidgetRepository
	commander        domain.DeviceCommander
	recorder         domain.Recorder
}

func NewCommandUsecase(dr domain.DeviceRepository, wr domain.WidgetRepository, dc domain.DeviceCommander, rd domain.Recorder) CommandUsecase {
	return &commandUsecase{
		deviceRepo:       dr,
		widgetRepository: wr,
		commander:        dc,
		recorder:         rd,
	}
}

func (u *commandUsecase) SendCommand(cmd *domain.DeviceCommand, widgetID uint, correlationID string) error {
	device, err := u.deviceRepo.FindByWidgetID(widgetID)
	if err != nil {
		return err
	}

	resp, err := u.commander.RequestCommand(device.DeviceID, cmd, correlationID)
	if err != nil {
		return err
	}

	if resp.Status != "success" {
		return err
	}
	err = u.recorder.RecordLog(&domain.Log{
		WidgetID:  widgetID,
		EventType: "command",
		Value:     cmd.Value,
	})
	if err != nil {
		return err
	}

	return u.widgetRepository.UpdateValue(widgetID, cmd.Value)
}
