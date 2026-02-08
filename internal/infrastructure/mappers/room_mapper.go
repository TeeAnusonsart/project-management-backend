package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func DomainToRoomModel(r *domain.Room) *models.Room {
	return &models.Room{
		ID:   r.ID,
		Name: r.Name,
	}
}

func ModelToDomainRoom(m *models.Room) *domain.Room {
	return &domain.Room{
		ID:   m.ID,
		Name: m.Name,
	}
}
