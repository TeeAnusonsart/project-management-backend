package usecase

import "project-home-iot/internal/core/domain"

type RoomUsecase interface {
	CreateRoom(name string) error
	ListRooms() ([]*domain.Room, error)
	GetRoom(id uint) (*domain.Room, error)
	UpdateRoom(id uint, name string) error
	DeleteRoom(id uint) error

	AddDeviceToRoom(roomID uint, deviceID string) error
	ListDevicesInRoom(roomID uint) ([]*domain.DeviceSummary, error)
}

type roomUsecase struct {
	roomRepo domain.RoomRepository
}

func NewRoomUsecase(r domain.RoomRepository) RoomUsecase {
	return &roomUsecase{roomRepo: r}
}

func (u *roomUsecase) CreateRoom(name string) error {
	room := &domain.Room{Name: name}
	return u.roomRepo.Create(room)
}

func (u *roomUsecase) ListRooms() ([]*domain.Room, error) {
	return u.roomRepo.FindAll()
}

func (u *roomUsecase) GetRoom(id uint) (*domain.Room, error) {
	return u.roomRepo.FindByID(id)
}

func (u *roomUsecase) UpdateRoom(id uint, name string) error {
	room, err := u.roomRepo.FindByID(id)
	if err != nil {
		return err
	}
	room.Name = name
	return u.roomRepo.Update(room)
}

func (u *roomUsecase) DeleteRoom(id uint) error {
	return u.roomRepo.Delete(id)
}

func (u *roomUsecase) AddDeviceToRoom(roomID uint, deviceID string) error {
	return u.roomRepo.AddDevice(roomID, deviceID)
}

func (u *roomUsecase) ListDevicesInRoom(roomID uint) ([]*domain.DeviceSummary, error) {
	return u.roomRepo.ListDeviceSummaries(roomID)
}
