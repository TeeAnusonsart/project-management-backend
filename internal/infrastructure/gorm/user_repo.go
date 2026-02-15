package gorm

import (
	"project-home-iot/internal/core/domain"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string
	Name        string
	Password    string
	Email       string
	ProfilePath string
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
	user.ID = model.ID
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

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var model User
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return ModelToDomainUser(&model), nil
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

func (r *UserRepository) UpdateProfilePath(id uint, path string) error {
	return r.db.Model(&User{}).
		Where("id = ?", id).
		Update("profile_path", path).Error
}


func DomainToUserModel(u *domain.User) *User {
	return &User{
		Model:       gorm.Model{ID: u.ID},
		Username:    u.Username,
		Name:        u.Name,
		Password:    u.Password,
		Email:       u.Email,
		ProfilePath: u.ProfilePath,
		Role:        string(u.Role),
	}
}

func ModelToDomainUser(m *User) *domain.User {
	return &domain.User{
		ID:          m.ID,
		Username:    m.Username,
		Name:        m.Name,
		Password:    m.Password,
		Email:       m.Email,
		ProfilePath: m.ProfilePath,
		Role:        domain.UserRole(m.Role),
	}
}
