package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mappers"
	"gorm.io/gorm"
)

type WidgetRepository struct {
	db *gorm.DB
}

func NewWidgetRepository(db *gorm.DB) *WidgetRepository {
	return &WidgetRepository{db: db}
}

func (r *WidgetRepository) CreateWidget(widget *domain.Widget) error {
	model := mappers.WidgetDomainToModel(widget)
	return r.db.Create(model).Error
}

func (r *WidgetRepository) UpdateValue(widgetID uint, value uint) error {
	return r.db.Model(&domain.Widget{}).
		Where("id = ?", widgetID).
		Update("value", value).Error
}

func (r *WidgetRepository) FindAll() ([]*domain.Widget, error) {
	var widgetModels []models.Widget

	err := r.db.
		Preload("Device").
		Preload("Capability").
		Where("widget_status = ?", "include").
		Order("widget_order asc").
		Find(&widgetModels).Error

	if err != nil {
		return nil, err
	}

	var result []*domain.Widget
	for _, m := range widgetModels {
		result = append(result, mappers.WidgetModelToDomain(&m))
	}

	return result, nil
}


func (r *WidgetRepository) FindByID(id uint) (*domain.Widget, error) {
	var widgetModel models.Widget

	err := r.db.
		Preload("Device").
		Preload("Capability").
		First(&widgetModel, id).Error

	if err != nil {
		return nil, err
	}

	return mappers.WidgetModelToDomain(&widgetModel), nil
}

func (r *WidgetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Widget{}, id).Error
}

func (r *WidgetRepository) Update(widget *domain.Widget) error {
	model := mappers.WidgetDomainToModel(widget)

	return r.db.Model(&models.Widget{}).
		Where("id = ?", widget.ID).
		Updates(model).Error
}

func (r *WidgetRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Widget{}).
		Where("id = ?", id).
		Update("widget_status", status).Error
}

func (r *WidgetRepository) ChangeOrder(widgetOrders []uint) error {
	for index, widgetID := range widgetOrders {
		if err := r.db.Model(&models.Widget{}).
			Where("id = ?", widgetID).
			Update("widget_order", index+1).Error; err != nil {
			return err
		}
	}
	return nil
}

