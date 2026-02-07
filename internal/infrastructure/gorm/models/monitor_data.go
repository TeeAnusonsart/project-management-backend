package models

import (
	"gorm.io/gorm"
)

// Device ตารางหลักของอุปกรณ์
type MonitorData struct {
	gorm.Model
	Value    uint
	WidgetID uint
}
