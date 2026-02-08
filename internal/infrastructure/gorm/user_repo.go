package gorm

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/mappers"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	model := mappers.DomainToUserModel(user)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	return nil
}

func (r *UserRepository) FindAll() ([]*domain.User, error) {
	var models []models.User
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var result []*domain.User
	for _, m := range models {
		result = append(result, mappers.ModelToDomainUser(&m))
	}
	return result, nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var model models.User
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return mappers.ModelToDomainUser(&model), nil
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

func (r *UserRepository) UpdateProfilePath(id uint, path string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("profile_path", path).Error
}

