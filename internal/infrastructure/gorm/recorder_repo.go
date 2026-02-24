package gorm

import (
	"project-home-iot/internal/core/domain"
	"gorm.io/gorm"
	// "project-home-iot/internal/infrastructure/gorm/models"
	// "project-home-iot/internal/infrastructure/mappers"
)

type Log struct {
	gorm.Model
	Value     string
	EventType string

	Actor string 
	User User	`gorm:"foreignKey:Actor;references:Email"`

	WidgetID  uint       `gorm:"not null"`
	Widget    Widget `gorm:"foreignKey:WidgetID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RecorderRepository struct {
	db *gorm.DB
}

func NewRecorderRepository(db *gorm.DB) *RecorderRepository {
	return &RecorderRepository{db: db}
}



func (r *RecorderRepository) RecordLog(log *domain.Log) error {

	model := LogDomainToModel(log)

    if err := r.db.Create(model).Error; err != nil {
        return err
    }

    log.ID = model.ID
    return nil
}

func (r *RecorderRepository) GetLogByWidgetID(widgetID uint) ([]*domain.Log, error) {
	var result []*domain.Log

	err := r.db.
		Model(&Log{}).
		Select("value, event_type,actor,created_at").
		Where("widget_id = ?", widgetID).
		Scan(&result).Error

	return result, err

}

func LogDomainToModel(d *domain.Log) *Log {
	actor := d.Actor
    if actor == "" {
        actor = "Sensor"
    }
	return &Log{
		WidgetID:  d.WidgetID,
		Value:     d.Value,
		Actor: actor,
		EventType: d.EventType,
	}
}

func LogModelToDomain(m *Log) *domain.Log {
	return &domain.Log{
		ID:        m.ID,
		// WidgetID:  m.WidgetID,
		Actor: m.Actor,
		Value:     m.Value,
		EventType: m.EventType,
	}
}