package models

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type Widget struct {
	gorm.Model
	WidgetStatus string
	Value         uint
	WidgetOrder   uint `gorm:"column:widget_order"`

	CapabilityID  uint       `gorm:"not null"`
	Capability    Capability `gorm:"foreignKey:CapabilityID;references:ID"`
	
	DeviceID     string       `gorm:"not null"`
    Device       Device     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
