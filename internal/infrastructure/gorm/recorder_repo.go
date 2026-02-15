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

	WidgetID  uint       `gorm:"not null"`
	Widget    Widget `gorm:"foreignKey:WidgetID;references:ID"`
}

type RecorderRepository struct {
	db *gorm.DB
}

func NewRecorderRepository(db *gorm.DB) *RecorderRepository {
	return &RecorderRepository{db: db}
}

// func (r *RecorderRepository) DataReceive(log *domain.Log) error {
// 	model := mappers.LogDomainToModel(log)

// 	if err := r.db.Create(model).Error; err != nil {
// 		return err
// 	}

// 	// sync ID กลับไป domain (optional แต่ดี)
// 	log.ID = model.ID
// 	return nil
// }

func (r *RecorderRepository) RecordLog(log *domain.Log) error {

	// log := &domain.Log{
	// 	WidgetID:  widgetId,
	// 	EventType: eventType,
	// 	Value:     value,
	// 	ActorType: "system", // หรือ inject จาก context
	// }

	// return r.DataReceive(log)
	model := LogDomainToModel(log)
    if err := r.db.Create(model).Error; err != nil {
        return err
    }

    log.ID = model.ID
    return nil
}

func LogDomainToModel(d *domain.Log) *Log {
	return &Log{
		WidgetID:  d.WidgetID,
		Value:     d.Value,
		EventType: d.EventType,
	}
}

func LogModelToDomain(m *Log) *domain.Log {
	return &domain.Log{
		ID:        m.ID,
		WidgetID:  m.WidgetID,
		Value:     m.Value,
		EventType: m.EventType,
	}
}