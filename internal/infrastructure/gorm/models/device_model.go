package models

import "time"

type Device struct {
	DeviceID string `gorm:"primaryKey;autoIncrement:false"`

	DeviceName string
	DeviceType string
	// Topic      string

	RoomID *uint
	Room   *Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LastHeartbeat time.Time

	Widgets []Widget `gorm:"foreignKey:DeviceID"`
}
