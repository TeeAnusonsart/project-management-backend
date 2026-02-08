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
	// Implementation for receiving data
	return nil
}

func (r *RecorderRepository) RecordLog(widgetId uint,eventType string,value uint) error {
	log := &domain.Log{
		WidgetID:  widgetId,
		EventTpye: eventType,
		Value:     value,
		ActorType: "",
	}
	model := mappers.LogDomainToModel(log)
	if err := r.db.Create(model).Error; err != nil {
        return err
    }
	// Implementation for recording log
	return nil
}