package models


type Device struct {
	ID uint `gorm:"primaryKey;autoIncrement:false"`

	DeviceName string
	DeviceType string
	Topic      string

	RoomID *uint
	Room   *Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Widgets []Widget `gorm:"foreignKey:DeviceID"`
}
