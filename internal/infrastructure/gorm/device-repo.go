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
