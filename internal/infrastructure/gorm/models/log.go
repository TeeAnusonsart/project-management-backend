package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	Value     uint
	EventType string

	WidgetID  uint       `gorm:"not null"`
	Widget    Widget `gorm:"foreignKey:WidgetID;references:ID"`
}
