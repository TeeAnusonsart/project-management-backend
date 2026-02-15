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
		Joins("join widgets on widgets.device_id = devices.device_id").
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
		Select("device_id,last_heartbeat device_name, device_type").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}


func (r *DeviceRepository) GetByID(id string) (*domain.Device, error) {
	var model models.Device

	if err := r.db.
		Where("device_id = ?", id).
		First(&model).Error; err != nil {
		return nil, err
	}

	return mappers.ModelToDomainDevice(&model), nil
}

func (r *DeviceRepository) GetSummaryByID(id string) (*domain.DeviceSummary, error) {
	var result domain.DeviceSummary

	err := r.db.
		Table("devices").
		Select("device_id, last_heartbeat, device_name, device_type").
		Where("device_id = ?", id).
		First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}


func (r *DeviceRepository) UpdateName(id string, name string) error {
	return r.db.Model(&models.Device{}).
		Where("device_id = ?", id).
		Update("device_name", name).Error
}

func (r *DeviceRepository) Pair(id string, deviceKey string) error {
	return r.db.Model(&models.Device{}).
		Where("device_id = ?", id).
		Update("device_key", deviceKey).Error
}

func (r *DeviceRepository) Unpair(id string) error {
	return r.db.Model(&models.Device{}).
		Where("device_id = ?", id).
		Update("device_key", "").Error
}

func (r *DeviceRepository) UpdateHeartbeat(deviceID string) error {
	return r.db.Model(&models.Device{}).
		Where("device_id = ?", deviceID).
		Update("last_heartbeat", gorm.Expr("NOW()")).Error
}

