package gorm

import (
	"project-home-iot/internal/core/domain"

	"gorm.io/gorm"
)

type Room struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"column:room_name"`

	Devices []Device `gorm:"foreignKey:RoomID"`
}


type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *domain.Room) error {
	model := DomainToRoomModel(room)

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	room.ID = model.ID
	return nil
}

func (r *RoomRepository) FindAll() ([]*domain.Room, error) {
	var models []Room
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var rooms []*domain.Room
	for _, m := range models {
		rooms = append(rooms, ModelToDomainRoom(&m))
	}

	return rooms, nil
}

func (r *RoomRepository) FindByID(id uint) (*domain.Room, error) {
	var model Room
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return ModelToDomainRoom(&model), nil
}

func (r *RoomRepository) Update(room *domain.Room) error {
	return r.db.Model(&Room{}).
		Where("id = ?", room.ID).
		Update("room_name", room.Name).Error
}

func (r *RoomRepository) Delete(id uint) error {
	if err := r.db.Model(&Device{}).
		Where("room_id = ?", id).
		Update("room_id", nil).Error; err != nil {
		return err
	}

	return r.db.Delete(&Room{}, id).Error
}


func (r *RoomRepository) AddDevice(roomID uint, deviceID uint) error {
	return r.db.Model(&Device{}).
		Where("device_id = ?", deviceID).
		Update("room_id", roomID).Error
}

func (r *RoomRepository) ListDeviceSummaries(roomID uint) ([]*domain.DeviceSummary, error) {
	var result []*domain.DeviceSummary

	err := r.db.
		Table("devices").
		Select("device_id, device_name, device_type").
		Where("room_id = ?", roomID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func DomainToRoomModel(r *domain.Room) *Room {
	return &Room{
		ID:   r.ID,
		Name: r.Name,
	}
}

func ModelToDomainRoom(m *Room) *domain.Room {
	return &domain.Room{
		ID:   m.ID,
		Name: m.Name,
	}
}



