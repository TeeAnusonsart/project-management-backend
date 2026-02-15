package usecase

import (
	"project-home-iot/internal/core/domain"
)

type RecordLogUsecase interface {
	// Execute(deviceId string, capabilityId uint, eventType string, value uint) error
	Execute(deviceId string, capabilityType string,controlType string, eventType string, value uint) error
}

type recordLogUsecase struct {
	widgetRepo domain.WidgetRepository
	capabilityRepo domain.CapabilityRepository
	recorder domain.Recorder
}

func NewRecordLogUsecase(r domain.Recorder,wr domain.WidgetRepository,cr domain.CapabilityRepository) RecordLogUsecase {
	return &recordLogUsecase{recorder: r, widgetRepo: wr,capabilityRepo: cr}
}

func (uc *recordLogUsecase) Execute(deviceId string, capabilityType string,controlType string, eventType string, value uint) error {
	capability,err := uc.capabilityRepo.FindByTypeAndControl(capabilityType, controlType)
	if err != nil {
		return err
	}
	widget := uc.widgetRepo.GetWidgetIdByDeviceAndCapability(deviceId, capability.ID)
	logEntry := &domain.Log{
        WidgetID:  widget.ID,
        EventType: eventType,
        Value:     value,
    }
    return uc.recorder.RecordLog(logEntry)
	// return uc.recorder.RecordLog(widget.ID, eventType, value)
}
