package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/mqtt/dtos"
)

func DeviceCommandDomainToPayload(d *domain.DeviceCommand) *dtos.DeviceCommand {
	return &dtos.DeviceCommand{
		CapabilityID: d.CapabilityID,
		Value:        d.Value,
		ReplyTopic:   d.ReplyTopic,
	}
}



