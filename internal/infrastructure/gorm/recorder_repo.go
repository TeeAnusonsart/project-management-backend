package gorm
import (
	"gorm.io/gorm"
	"project-home-iot/internal/core/domain"
	// "project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mappers"
)

type RecorderRepository struct {
	db *gorm.DB
}

func NewRecorderRepository(db *gorm.DB) *RecorderRepository {
	return &RecorderRepository{db: db}
}

func (r *RecorderRepository) DataReceive(log *domain.Log) error {
	model := mappers.LogDomainToModel(log)

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	// sync ID กลับไป domain (optional แต่ดี)
	log.ID = model.ID
	return nil
}

func (r *RecorderRepository) RecordLog(
	widgetId uint,
	eventType string,
	value uint,
) error {

	log := &domain.Log{
		WidgetID:  widgetId,
		EventType: eventType,
		Value:     value,
		ActorType: "system", // หรือ inject จาก context
	}

	return r.DataReceive(log)
}