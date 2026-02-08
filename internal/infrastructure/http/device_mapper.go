package http

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/http/dtos"
)

func toDeviceResponse(d *domain.DeviceSummary) dtos.DeviceResponse {
	return dtos.DeviceResponse{
		ID:         d.ID,
		DeviceName: d.DeviceName,
		DeviceType: d.DeviceType,
	}
}
