package domain

import "fmt"

type Device struct {
	ID uint

	DeviceName string
	DeviceType string
	Topic      string

	Widgets []Widget
}

type DeviceSummary struct {
	ID   uint
	DeviceName string
	DeviceType string
}

type DeviceCommand struct {
	CapabilityType        string
	ControlType string
	Value          uint
	ReplyTopic     string
}

type CommandResponse struct {
	Status        string
}

func (d *Device) Validate() error {
	if d.ID == 0 {
		return fmt.Errorf("device id is required")
	}
	if d.DeviceName == "" {
		return fmt.Errorf("device name is required")
	}
	if d.DeviceType == "" {
		return fmt.Errorf("device type is required")
	}
	if d.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	return nil
}



type DeviceRepository interface {
	
	CreateDevice(device *Device) error
	FindByWidgetID(widgetID uint) (*Device, error)

	GetAllSummaries() ([]*DeviceSummary, error)
	GetSummaryByID(id uint) (*DeviceSummary, error)

	GetByID(id uint) (*Device, error)

	UpdateName(id uint, name string) error
	Pair(id uint, deviceKey string) error
	Unpair(id uint) error
}

type PairCommander interface {
	RequestPair(deviceID uint, deviceKey string) error
	Subscribe(topic string) error
}

type DeviceCommander interface {
	RequestCommand(deviceId uint, cmd *DeviceCommand,correlationID string) (*CommandResponse, error)
}

