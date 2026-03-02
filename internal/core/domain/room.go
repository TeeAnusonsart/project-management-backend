package domain

type Room struct {
	ID      uint
	Name    string
	Devices []DeviceSummary
}

type RoomRepository interface {
	Create(room *Room) error
	FindAll() ([]*Room, error)
	FindByID(id uint) (*Room, error)
	FindByRoomName(name string) (*Room, error)
	Update(room *Room) error
	Delete(id uint) error

	AddDevice(roomID uint, deviceID string) error
	ListDeviceSummaries(roomID uint) ([]*DeviceSummary, error)
}
