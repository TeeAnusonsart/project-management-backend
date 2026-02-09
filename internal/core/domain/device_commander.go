package domain

type DeviceCommand struct {
	WidgetID       uint
	CapabilityID        string
	Value          uint
	ReplyTopic     string
	CorrelationID  string
}

type CommandResponse struct {
	CorrelationID string
	Status        string
	Value          uint
	Message       string
}

type DeviceCommander interface {
	RequestCommand(topic string, cmd *DeviceCommand) (*CommandResponse, error)
}