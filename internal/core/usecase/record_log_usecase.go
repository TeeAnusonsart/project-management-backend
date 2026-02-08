package usecase

import "project-home-iot/internal/core/domain"

type RecordLogUsecase struct {
	widgetRepo domain.WidgetRepository
	recorder domain.Recorder
}

func NewRecordLogUsecase(r domain.Recorder,wr domain.WidgetRepository) *RecordLogUsecase {
	return &RecordLogUsecase{recorder: r, widgetRepo: wr}
}

func (uc *RecordLogUsecase) Execute(
	deviceId uint,
	capabilityId uint,
	eventType string,
	value uint,
) error {
	widget := uc.widgetRepo.GetWidgetIdByDeviceAndCapability(deviceId, capabilityId)
	return uc.recorder.RecordLog(widget.ID, eventType, value)
}
