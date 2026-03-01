package gorm

import (
	"project-home-iot/internal/core/domain"
	"gorm.io/gorm"
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

// func (r *UserRepository) FindByUID(id uint) (*domain.User, error) {
// 	var model User
// 	if err := r.db.First(&model, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return ModelToDomainUser(&model), nil
// }

func (r *UserRepository) Delete(email string) error {
	return r.db.Where("email = ?", email).Delete(&User{}).Error
}

// func (r *UserRepository) UpdatePassword(id uint, hashedPassword string) error {
// 	return r.db.Model(&User{}).
// 		Where("id = ?", id).
// 		Update("password", hashedPassword).Error
// }

// func (r *UserRepository) UpdateProfilePath(id uint, path string) error {
// 	return r.db.Model(&User{}).
// 		Where("id = ?", id).
// 		Update("profile_path", path).Error
// }

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
    var model User
    if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
        return nil, err
    }
    return ModelToDomainUser(&model), nil
}


// func (r *UserRepository) UpdateUID(email string, uid string) error {
//     return r.db.Model(&User{}).Where("id = ?", id).Update("uid", uid).Error
// }


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
