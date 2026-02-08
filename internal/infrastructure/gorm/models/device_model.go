package models


type Device struct {
	ID uint `gorm:"primaryKey"`

	// DeviceID string `gorm:"uniqueIndex;not null"`

	DeviceName string
	DeviceKey  string
	DeviceType string
	Topic      string
	RoomID *uint
	Room   *Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`


	Widgets    []Widget `gorm:"foreignKey:DeviceID"`

}
