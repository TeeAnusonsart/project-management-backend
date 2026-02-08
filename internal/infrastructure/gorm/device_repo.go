package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/mappers"
	"project-home-iot/internal/infrastructure/gorm/models"
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

	return r.db.Create(model).Error
}

func (r *DeviceRepository) FindByWidgetID(widgetID uint) (*domain.Device, error) {
	var model models.Device

	err := r.db.
		Joins("join widgets on widgets.device_id = devices.id").
		Where("widgets.id = ?", widgetID).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return mappers.ModelToDomainDevice(&model), nil
}

// func (r *DeviceRepository) CreateMonitorData(data *domain.MonitorData) error {
// 	return r.db.Create(data).Error
// }

func (r *DeviceRepository) GetAllSummaries() ([]*domain.DeviceSummary, error) {
	var result []*domain.DeviceSummary

	err := r.db.
		Table("devices").
		Select("id, device_name, device_type").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}


func (r *DeviceRepository) GetByID(id uint) (*domain.Device, error) {
	var model models.Device
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return mappers.ModelToDomainDevice(&model), nil
}

func (r *DeviceRepository) GetSummaryByID(id uint) (*domain.DeviceSummary, error) {
	var result domain.DeviceSummary

	err := r.db.
		Table("devices").
		Select("id, device_name, device_type").
		Where("id = ?", id).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
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

