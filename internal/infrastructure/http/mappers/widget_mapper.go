package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/http/dtos"
)

func ToWidgetResponse(w *domain.Widget) dtos.WidgetResponse {

	response := dtos.WidgetResponse{
		WidgetID:     w.ID,
		WidgetOrder:  w.WidgetOrder,
		WidgetStatus: w.WidgetStatus,
		Value:        w.Value,
	}

	if w.Device != nil {
		response.Device = dtos.DeviceDTO{
			DeviceID:   w.Device.DeviceID,
			DeviceLastHeartbeat: w.Device.LastHeartbeat,
			DeviceName: w.Device.DeviceName,
			DeviceType: w.Device.DeviceType,
		}
	}

	if w.Capability != nil {
		response.Capability = dtos.CapabilityDTO{
			CapabilityID:   w.Capability.ID,
			CapabilityType: w.Capability.CapabilityType,
			ControlType:    w.Capability.ControlType,
		}
	}

	return response
}
