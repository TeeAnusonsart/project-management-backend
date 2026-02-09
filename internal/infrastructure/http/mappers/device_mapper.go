package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/http/dtos"
)

func ToDeviceResponse(d *domain.DeviceSummary) dtos.DeviceResponse {
	return dtos.DeviceResponse{
		ID:         d.ID,
		DeviceName: d.DeviceName,
		DeviceType: d.DeviceType,
	}
}
