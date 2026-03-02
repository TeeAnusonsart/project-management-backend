package usecase

import (
	"project-home-iot/internal/core/domain"
	"errors"
)

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
	deviceRepo domain.DeviceRepository
}

func NewRoomUsecase(r domain.RoomRepository,d domain.DeviceRepository) RoomUsecase {
	return &roomUsecase{roomRepo: r,deviceRepo: d}
}

func (u *roomUsecase) CreateRoom(name string) error {
	_, err := u.roomRepo.FindByRoomName(name)

	if err == nil {
		return domain.ErrRoomExist
	}

	if !errors.Is(err, domain.ErrRoomNotFound) {
		return err
	}

	room := &domain.Room{Name: name}
	return u.roomRepo.Create(room)
}

func (u *roomUsecase) ListRooms() ([]*domain.Room, error) {
	
	return u.roomRepo.FindAll()
}

func (u *roomUsecase) GetRoom(id uint) (*domain.Room, error) {
	_, err := u.roomRepo.FindByID(id)
    if err != nil {
        return nil,domain.ErrRoomNotFound 
    }
	return u.roomRepo.FindByID(id)
}

func (u *roomUsecase) UpdateRoom(id uint, name string) error {
	room, err := u.roomRepo.FindByID(id)
	if err != nil {
		return domain.ErrRoomNotFound
	}
	room.Name = name
	return u.roomRepo.Update(room)
}

func (u *roomUsecase) DeleteRoom(id uint) error {
	_, err := u.roomRepo.FindByID(id)
    if err != nil {
        return domain.ErrRoomNotFound 
    }
	return u.roomRepo.Delete(id)
}

func (u *roomUsecase) AddDeviceToRoom(roomID uint, deviceID string) error {
	_, err := u.roomRepo.FindByID(roomID)
    if err != nil {
        return domain.ErrRoomNotFound 
    }
	_, err = u.deviceRepo.GetByID(deviceID)
	if err != nil {
        return domain.ErrDeviceNotFound 
    }
	return u.roomRepo.AddDevice(roomID, deviceID)
}

func (u *roomUsecase) ListDevicesInRoom(roomID uint) ([]*domain.DeviceSummary, error) {
	_, err := u.roomRepo.FindByID(roomID)
    if err != nil {
        return nil,domain.ErrRoomNotFound 
    }
	
	return u.roomRepo.ListDeviceSummaries(roomID)
}
