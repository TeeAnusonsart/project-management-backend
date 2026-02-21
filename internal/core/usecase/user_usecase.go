package usecase

import (
	"errors"

	"project-home-iot/internal/core/domain"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	ListUsers() ([]*domain.User, error)
	GetUser(id uint) (*domain.User, error)
	DeleteUser(id uint) error
	UploadProfile(id uint, filePath string) error
}

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) UserUsecase {
	return &userUsecase{userRepo: ur}
}

func (u *userUsecase) CreateUser(user *domain.User) error {

	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	if !user.Role.IsValid() {
		return errors.New("invalid role")
	}

	return u.userRepo.Create(user)
}


func (u *userUsecase) ListUsers() ([]*domain.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) GetUser(id uint) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) DeleteUser(id uint) error {
	return u.userRepo.Delete(id)
}


func (u *userUsecase) UploadProfile(id uint, filePath string) error {
	return u.userRepo.UpdateProfilePath(id, filePath)
}
