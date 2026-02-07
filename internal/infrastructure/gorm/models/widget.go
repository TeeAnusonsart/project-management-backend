package models

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type Widget struct {
	gorm.Model
	Widget_status string
	Value         uint
	CapabilityID  uint       `gorm:"not null"`
	Capability    Capability `gorm:"foreignKey:CapabilityID;references:ID"`
	DeviceID     uint       `gorm:"not null"`
    Device       Device     `gorm:"foreignKey:DeviceID;references:ID"`
}
