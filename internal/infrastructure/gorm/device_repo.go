package gorm

import (
	"project-home-iot/internal/core/domain"
	// "project-home-iot/internal/infrastructure/gorm/models"
	"time"

	"gorm.io/gorm"
)

type Device struct {
	DeviceID string `gorm:"primaryKey;autoIncrement:false"`

	DeviceName string
	DeviceType string

	RoomID        *uint
	Room          *Room `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LastHeartbeat time.Time

	Widgets []Widget `gorm:"foreignKey:DeviceID"`
}

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) CreateDevice(device *domain.Device) error {
	model := DomainToDeviceModel(device)

	return r.db.Create(model).Error
}

func (r *DeviceRepository) FindByWidgetID(widgetID uint) (*domain.Device, error) {
	var model Device

	err := r.db.
		Joins("join widgets on widgets.device_id = devices.device_id").
		Where("widgets.id = ?", widgetID).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return ModelToDomainDevice(&model), nil
}

func (r *DeviceRepository) GetAllSummaries() ([]*domain.DeviceSummary, error) {
	var result []*domain.DeviceSummary

	err := r.db.
		Table("devices").
		Select("device_id,last_heartbeat, device_name, device_type").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *DeviceRepository) GetUnpairDevice() ([]*domain.DeviceSummary, error) {
	var result []*domain.DeviceSummary

	err := r.db.
		Model(&Device{}).
		Select("devices.*").
		Joins("LEFT JOIN widgets ON widgets.device_id = devices.device_id").
		Where("widgets.device_id IS NULL").
		Scan(&result).Error
	return result, err
}

func (r *DeviceRepository) GetPairedDevice() ([]*domain.DeviceSummary, error) {
	var result []*domain.DeviceSummary

	err := r.db.
		Model(&Device{}).
		Where("EXISTS (?)",
			r.db.
				Select("1").
				Table("widgets").
				Where("widgets.device_id = devices.device_id"),
		).
		Scan(&result).Error
	return result, err
}

func (r *DeviceRepository) GetByID(id string) (*domain.Device, error) {
	var model Device

	if err := r.db.
		Where("device_id = ?", id).
		First(&model).Error; err != nil {
		return nil, err
	}

	return ModelToDomainDevice(&model), nil
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
	return r.db.Model(&Device{}).
		Where("device_id = ?", id).
		Update("device_name", name).Error
}

func (r *DeviceRepository) Pair(id string, deviceKey string) error {
	return r.db.Model(&Device{}).
		Where("device_id = ?", id).
		Update("device_key", deviceKey).Error
}

func (r *DeviceRepository) Unpair(id string) error {
	return r.db.Model(&Device{}).
		Where("device_id = ?", id).
		Update("device_key", "").Error
}

func (r *DeviceRepository) UpdateHeartbeat(deviceID string) error {
	return r.db.Model(&Device{}).
		Where("device_id = ?", deviceID).
		Update("last_heartbeat", gorm.Expr("NOW()")).Error
}

func DomainToDeviceModel(d *domain.Device) *Device {
	return &Device{
		DeviceID:   d.DeviceID,
		DeviceName: d.DeviceName,
		DeviceType: d.DeviceType,
	}
}

func ModelToDomainDevice(m *Device) *domain.Device {
	return &domain.Device{
		DeviceID:   m.DeviceID,
		DeviceName: m.DeviceName,
		DeviceType: m.DeviceType,
	}
}
