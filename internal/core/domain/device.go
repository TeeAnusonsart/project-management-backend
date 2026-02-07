package domain

type Device struct {
	ID uint

	DeviceID   string
	DeviceName string
	DeviceKey  string
	DeviceType string
	Topic      string

	Widgets []Widget
}
