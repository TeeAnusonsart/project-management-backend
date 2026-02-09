package domain

type DeviceCommand struct {
	CapabilityID        string
	Value          uint
	ReplyTopic     string
}

type CommandResponse struct {
	Status        string
}

type DeviceCommander interface {
	RequestCommand(topic string, cmd *DeviceCommand,correlationID string) (*CommandResponse, error)
}