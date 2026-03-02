package gorm

import (
	"project-home-iot/internal/core/domain"
	"gorm.io/gorm"
	"errors"
)

type User struct {
	Email       string `gorm:"primaryKey;autoIncrement:false"`
	Name string
	Role        string
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	model := DomainToUserModel(user)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	user.Email = model.Email
	return nil
}

func (r *UserRepository) FindAll() ([]*domain.User, error) {
	var models []User
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	var result []*domain.User
	for _, m := range models {
		result = append(result, ModelToDomainUser(&m))
	}
	return result, nil
}



func (r *UserRepository) Delete(email string) error {
	return r.db.Where("email = ?", email).Delete(&User{}).Error
}



func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var model User

	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}

	user := ModelToDomainUser(&model)
	return user, nil
}




func DomainToUserModel(u *domain.User) *User {
	return &User{
		Email:       u.Email,
		Name: u.Name,
		Role:        string(u.Role),
	}
}

func ModelToDomainUser(m *User) *domain.User {
	return &domain.User{
		Email:       m.Email,
		Name: m.Name,
		Role:        domain.UserRole(m.Role),
	}
}
