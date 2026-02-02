package domain

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type Capability struct {
	gorm.Model
	CapabilityType string
	Widgets        []Widget
}
