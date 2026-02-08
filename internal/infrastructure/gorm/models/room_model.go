package models

type Room struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"column:room_name"`

	Devices []Device `gorm:"foreignKey:RoomID"`
}
