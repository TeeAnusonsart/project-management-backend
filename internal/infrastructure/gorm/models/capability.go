package models

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type Capability struct {
	gorm.Model
	CapabilityType string
	ControlType	string
	Widgets        []Widget `gorm:"foreignKey:CapabilityID"`
}
