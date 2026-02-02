package domain

import (
	"time"
)

type Device struct {
	ID uint

	DeviceID   string
	DeviceName string
	DeviceKey  string
	DeviceType string
	Topic      string

	LastOnline time.Time
	Widgets    []Widget
}
