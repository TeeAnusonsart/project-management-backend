package usecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"project-home-iot/internal/core/domain"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	ListUsers() ([]*domain.User, error)
	GetUser(id uint) (*domain.User, error)
	DeleteUser(id uint) error
	ChangePassword(id uint, oldPassword, newPassword string) error
	UploadProfile(id uint, filePath string) error
}

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(ur domain.UserRepository) UserUsecase {
	return &userUsecase{userRepo: ur}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}

	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	if !user.Role.IsValid() {
		return errors.New("invalid role")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
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

func (u *userUsecase) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(oldPassword),
	); err != nil {
		return errors.New("old password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.userRepo.UpdatePassword(id, string(hashed))
}

func (u *userUsecase) UploadProfile(id uint, filePath string) error {
	return u.userRepo.UpdateProfilePath(id, filePath)
}
