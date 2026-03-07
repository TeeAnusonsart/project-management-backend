package usecase

import (
	"project-home-iot/internal/core/domain"
)

type RecordLogUsecase interface {
	// Execute(deviceId string, capabilityId uint, eventType string, value uint) error
	Execute(deviceId string, capabilityType string,controlType string, value string, actor string, eventType string) error
}

type recordLogUsecase struct {
	widgetRepo domain.WidgetRepository
	capabilityRepo domain.CapabilityRepository
	recorder domain.Recorder
}

func NewRecordLogUsecase(r domain.Recorder,wr domain.WidgetRepository,cr domain.CapabilityRepository) RecordLogUsecase {
	return &recordLogUsecase{recorder: r, widgetRepo: wr,capabilityRepo: cr}
}

func (uc *recordLogUsecase) Execute(deviceId string, capabilityType string,controlType string, value string, actor string, eventType string) error {
	capability,err := uc.capabilityRepo.FindByTypeAndControl(capabilityType, controlType)
	if err != nil {
		return err
	}
	widget, err := uc.widgetRepo.GetWidgetIdByDeviceAndCapability(deviceId, capability.ID)
	if err != nil {
		return err
	}

	if widget == nil {
		return nil // ไม่มี widget ก็ไม่ต้อง log
	}
	logEntry := &domain.Log{
        WidgetID:  widget.ID,
        Value:     value,
		EventType: eventType,
		Actor: actor,
    }
	
	uc.widgetRepo.UpdateValue(widget.ID,value)
    return uc.recorder.RecordLog(logEntry)
	// return uc.recorder.RecordLog(widget.ID, eventType, value)
}
