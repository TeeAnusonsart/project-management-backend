package domain

type RoomRepository interface {
	Create(room *Room) error
	FindAll() ([]*Room, error)
	FindByID(id uint) (*Room, error)
	Update(room *Room) error
	Delete(id uint) error

	AddDevice(roomID uint, deviceID uint) error
	ListDevices(roomID uint) ([]*Device, error)
}
