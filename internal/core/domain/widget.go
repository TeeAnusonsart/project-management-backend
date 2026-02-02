package domain

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type Widget struct {
	gorm.Model
	Widget_status  string
	Value          uint
	CapabilityType string
	CapabilityID   uint
	DeviceID       uint
}
