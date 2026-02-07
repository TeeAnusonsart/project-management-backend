package models


type Device struct {
	ID uint `gorm:"primaryKey"`

	// DeviceID string `gorm:"uniqueIndex;not null"`

	DeviceName string
	DeviceKey  string
	DeviceType string
	Topic      string

	Widgets    []Widget `gorm:"foreignKey:DeviceID"`

}
