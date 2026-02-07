package gorm

import (
	"project-home-iot/internal/core/domain"
	// "project-home-iot/internal/infrastructure/gorm/models"
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
	return r.db.Model(&domain.Widget{}).Where("id = ?", widgetID).Update("value", value).Error
}

// func (r *WidgetRepository) FindByID(widgetID uint) (*domain.Widget, error) {
// 	var model models.Widget
// 	err := r.db.First(&model, widgetID).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return mappers.WidgetModelToDomain(&model), nil
// }
