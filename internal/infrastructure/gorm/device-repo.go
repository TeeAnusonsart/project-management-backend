package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/mappers"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) CreateDevice(device *domain.Device) error {
	model := mappers.DomainToDeviceModel(device)
	if err := r.db.Create(model).Error; err != nil {
        return err
    }
	device.ID = model.ID
	return nil
}

func (r *DeviceRepository) FindByWidgetID(widgetID uint) (*domain.Device, error) {
	model := domain.Device{}
	err := r.db.Table("devices").Select("devices.id as id, devices.device_name, devices.device_key, devices.device_type, devices.topic").
		Joins("join widgets on widgets.device_id = devices.id").
		Where("widgets.id = ?", widgetID).
		First(&model).Error
	if err != nil {
		return nil, err
	}
	device := &domain.Device{
		ID:         model.ID,
		// DeviceID:   model.DeviceID,
		DeviceName: model.DeviceName,
		DeviceKey:  model.DeviceKey,
		DeviceType: model.DeviceType,
		Topic:      model.Topic,
	}
	return device, nil
}

// func (r *DeviceRepository) CreateMonitorData(data *domain.MonitorData) error {
// 	return r.db.Create(data).Error
// }

func (r *DeviceRepository) GetAll() ([]*domain.Device, error) {
	var models []models.Device
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var devices []*domain.Device
	for _, m := range models {
		devices = append(devices, mappers.ModelToDomainDevice(&m))
	}
	return devices, nil
}

func (r *DeviceRepository) GetByID(id uint) (*domain.Device, error) {
	var model models.Device
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return mappers.ModelToDomainDevice(&model), nil
}

func (r *DeviceRepository) UpdateName(id uint, name string) error {
	return r.db.Model(&models.Device{}).
		Where("id = ?", id).
		Update("device_name", name).Error
}

func (r *DeviceRepository) Pair(id uint, deviceKey string) error {
	return r.db.Model(&models.Device{}).
		Where("id = ?", id).
		Update("device_key", deviceKey).Error
}

func (r *DeviceRepository) Unpair(id uint) error {
	return r.db.Model(&models.Device{}).
		Where("id = ?", id).
		Update("device_key", "").Error
}

