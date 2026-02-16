package gorm

import (
	"project-home-iot/internal/core/domain"
	"fmt"
	"gorm.io/gorm"
)

type Widget struct {
	gorm.Model
	WidgetStatus string
	Value         string
	WidgetOrder   uint `gorm:"column:widget_order"`

	CapabilityID  uint       `gorm:"not null"`
	Capability    Capability `gorm:"foreignKey:CapabilityID;references:ID"`
	
	DeviceID     string       `gorm:"not null"`
    Device       Device     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WidgetRepository struct {
	db *gorm.DB
}

func NewWidgetRepository(db *gorm.DB) *WidgetRepository {
	return &WidgetRepository{db: db}
}

func (r*WidgetRepository) GetWidgetIdByDeviceAndCapability(deviceId string, capabilityId uint) *domain.Widget {
	var widgetModel Widget
	r.db.Where("device_id = ? AND capability_id = ?", deviceId, capabilityId).
		First(&widgetModel)
	domainWidget := WidgetModelToDomain(&widgetModel)
	return domainWidget
}

func (r *WidgetRepository) CreateWidget(widget *domain.Widget) error {
	model := WidgetDomainToModel(widget)
	return r.db.Create(model).Error
}

func (r *WidgetRepository) UpdateValue(widgetID uint, value string) error {
	return r.db.Model(&Widget{}).
		Where("id = ?", widgetID).
		Update("value", value).Error
}

func (r *WidgetRepository) FindAll() ([]*domain.Widget, error) {
	var widgetModels []Widget

	err := r.db.
		Preload("Device").
		Preload("Capability").
		Order("widget_order asc").
		Find(&widgetModels).Error

	if err != nil {
		return nil, err
	}

	var result []*domain.Widget
	for _, m := range widgetModels {
		result = append(result, WidgetModelToDomain(&m))
	}

	return result, nil
}


func (r *WidgetRepository) FindByID(id uint) (*domain.Widget, error) {
	var widgetModel Widget

	err := r.db.
		Preload("Device").
		Preload("Capability").
		First(&widgetModel, id).Error

	if err != nil {
		return nil, err
	}

	return WidgetModelToDomain(&widgetModel), nil
}

func (r *WidgetRepository) FindByRoomID(roomID uint) ([]*domain.Widget, error) {
	var widgetModels []Widget

	err := r.db.
		Joins("JOIN devices ON devices.device_id = widgets.device_id").
		Where("devices.room_id = ?", roomID).
		Preload("Device").
		Preload("Capability").
		Order("widget_order asc").
		Find(&widgetModels).Error

	if err != nil {
		return nil, err
	}

	var result []*domain.Widget
	for _, m := range widgetModels {
		result = append(result, WidgetModelToDomain(&m))
	}

	return result, nil
}

func (r *WidgetRepository) Delete(id uint) error {
	return r.db.Delete(&Widget{}, id).Error
}

func (r *WidgetRepository) Update(widget *domain.Widget) error {
	model := WidgetDomainToModel(widget)

	return r.db.Model(&Widget{}).
		Where("id = ?", widget.ID).
		Updates(model).Error
}

func (r *WidgetRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&Widget{}).
		Where("id = ?", id).
		Update("widget_status", status).Error
}

func (r *WidgetRepository) ChangeOrder(roomID uint, widgetOrders []uint) error {
    for index, widgetID := range widgetOrders {

        var count int64
        err := r.db.
            Table("widgets").
            Joins("JOIN devices ON devices.device_id = widgets.device_id").
            Where("widgets.id = ?", widgetID).
            Where("devices.room_id = ?", roomID).
            Count(&count).Error

        if err != nil {
            return err
        }

        if count == 0 {
            return fmt.Errorf("widget %d not in room %d", widgetID, roomID)
        }

        err = r.db.Model(&Widget{}).
            Where("id = ?", widgetID).
            Update("widget_order", index).Error

        if err != nil {
            return err
        }
    }

    return nil
}

func (r *WidgetRepository) GetWidgetByStatus(status string) ([]*domain.Widget, error) {
	var widgetModels []Widget
	err := r.db.
		Where("widget_status = ?", status).
		Preload("Device").
		Preload("Capability").
		Order("widget_order asc").
		Find(&widgetModels).Error

	if err != nil {
		return nil, err
	}
	var result []*domain.Widget
	for _, m := range widgetModels {
		result = append(result, WidgetModelToDomain(&m))
	}
	return result, nil
}

func WidgetDomainToModel(d *domain.Widget) *Widget {
	return &Widget{
		Model: gorm.Model{ID: d.ID},
		DeviceID:      d.DeviceID,
		CapabilityID:  d.CapabilityID,
		Value:         d.Value,
		WidgetStatus: d.WidgetStatus,
		WidgetOrder:   d.WidgetOrder,
	}
}

func WidgetModelToDomain(m *Widget) *domain.Widget {

	widget := &domain.Widget{
		ID:            m.ID,
		DeviceID:      m.DeviceID,
		CapabilityID:  m.CapabilityID,
		Value:         m.Value,
		WidgetStatus:  m.WidgetStatus,
		WidgetOrder:   m.WidgetOrder,
	}

	if m.Device.DeviceID != "" {
		widget.Device = &domain.Device{
			DeviceID:         m.Device.DeviceID,
			LastHeartbeat: m.Device.LastHeartbeat,
			DeviceName: m.Device.DeviceName,
			DeviceType: m.Device.DeviceType,

		}
	}

	if m.Capability.ID != 0 {
		widget.Capability = &domain.Capability{
			ID:             m.Capability.ID,
			CapabilityType: m.Capability.CapabilityType,
			ControlType:    m.Capability.ControlType,
		}
	}

	return widget
}


