package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mappers"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *domain.Room) error {
	model := mappers.DomainToRoomModel(room)

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	room.ID = model.ID
	return nil
}

func (r *RoomRepository) FindAll() ([]*domain.Room, error) {
	var models []models.Room
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var rooms []*domain.Room
	for _, m := range models {
		rooms = append(rooms, mappers.ModelToDomainRoom(&m))
	}

	return rooms, nil
}

func (r *RoomRepository) FindByID(id uint) (*domain.Room, error) {
	var model models.Room
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return mappers.ModelToDomainRoom(&model), nil
}

func (r *RoomRepository) Update(room *domain.Room) error {
	return r.db.Model(&models.Room{}).
		Where("id = ?", room.ID).
		Update("room_name", room.Name).Error
}

func (r *RoomRepository) Delete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *RoomRepository) AddDevice(roomID uint, deviceID uint) error {
	return r.db.Model(&models.Device{}).
		Where("id = ?", deviceID).
		Update("room_id", roomID).Error
}

func (r *RoomRepository) ListDevices(roomID uint) ([]*domain.Device, error) {
	var deviceModels []models.Device
	if err := r.db.Where("room_id = ?", roomID).Find(&deviceModels).Error; err != nil {
		return nil, err
	}

	var result []*domain.Device
	for _, d := range deviceModels {
		result = append(result, mappers.ModelToDomainDevice(&d))
	}

	return result, nil
}

