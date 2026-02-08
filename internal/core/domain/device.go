package domain

import "fmt"

type Device struct {
	ID uint

	DeviceName string
	DeviceType string
	Topic      string

	Widgets []Widget
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

