package domain

import (
	"fmt"
	"time"
)

type Device struct {
	DeviceID string

	DeviceName    string
	DeviceType    string
	LastHeartbeat time.Time
	Widgets       []Widget
}

type DeviceSummary struct {
	DeviceID      string
	LastHeartbeat time.Time
	DeviceName    string
	DeviceType    string
}

type DeviceCommand struct {
	CapabilityType string
	ControlType    string
	Value          string
	ReplyTopic     string
}

type CommandResponse struct {
	Status string
}


func (d *Device) Validate() error {
	if d.DeviceID == "" {
		return fmt.Errorf("device id is required")
	}
	if d.DeviceName == "" {
		return fmt.Errorf("device name is required")
	}
	if d.DeviceType == "" {
		return fmt.Errorf("device type is required")
	}
	// if d.Topic == "" {
	// 	return fmt.Errorf("topic is required")
	// }
	return nil
}

type DeviceRepository interface {
	CreateDevice(device *Device) error
	FindByWidgetID(widgetID uint) (*Device, error)

	GetAllSummaries() ([]*DeviceSummary, error)
	GetSummaryByID(deviceID string) (*DeviceSummary, error)
	GetUnpairDevice() ([]*DeviceSummary, error)
	GetPairedDevice() ([]*DeviceSummary, error)

	GetByID(deviceID string) (*Device, error)
	UpdateHeartbeat(deviceID string) error
	UpdateName(deviceID string, name string) error
	Pair(deviceID string, deviceKey string) error
	Unpair(deviceID string) error
}

type PairCommander interface {
	RequestPair(deviceID string, deviceKey string) error
	Subscribe(deviceID string) error
}

type DeviceCommander interface {
	RequestCommand(deviceId string, cmd *DeviceCommand, correlationID string) (*CommandResponse, error)
}
