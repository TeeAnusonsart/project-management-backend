package models

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type WidgetModel struct {
	gorm.Model
	Widget_status string
	Value         uint
	CapabilityID  uint
	DeviceID      uint
}
