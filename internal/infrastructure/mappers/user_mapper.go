package mappers

import (
	"gorm.io/gorm"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
)

func DomainToUserModel(u *domain.User) *models.User {
	return &models.User{
		Model:       gorm.Model{ID: u.ID},
		Username:    u.Username,
		Name:        u.Name,
		Password:    u.Password,
		Email:       u.Email,
		ProfilePath: u.ProfilePath,
	}
}

func ModelToDomainUser(m *models.User) *domain.User {
	return &domain.User{
		ID:          m.ID,
		Username:    m.Username,
		Name:        m.Name,
		Password:    m.Password,
		Email:       m.Email,
		ProfilePath: m.ProfilePath,
	}
}
