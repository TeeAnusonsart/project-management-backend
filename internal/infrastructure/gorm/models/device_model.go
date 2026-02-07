package models


type DeviceModel struct {
	ID uint `gorm:"primaryKey"`

	DeviceID string `gorm:"uniqueIndex;not null"`

	DeviceName string
	DeviceKey  string
	DeviceType string
	Topic      string

}
