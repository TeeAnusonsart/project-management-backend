package usecase

import (
	"context"
	"errors"
	"project-home-iot/internal/core/domain"

	"firebase.google.com/go/v4/auth"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	ListUsers() ([]*domain.User, error)
	Login(ctx context.Context, idToken string) (*domain.User, error)
	GetUser(email string) (*domain.User, error)
	DeleteUser(email string) error
	UpdateUser(user *domain.User) error
}

type userUsecase struct {
	userRepo   domain.UserRepository
	authClient *auth.Client
}

func NewUserUsecase(ur domain.UserRepository, ac *auth.Client) UserUsecase {
	return &userUsecase{
		userRepo:   ur,
		authClient: ac,
	}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	exist, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}

	if exist != nil {
		return domain.ErrEmailAlreadyExist
	}

	if user.Role == "" {
		user.Role = domain.RoleUser
	}

	return u.userRepo.Create(user)
}

func (u *userUsecase) UpdateUser(user *domain.User) error {
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.New("user not found")
	}

	existingUser.Name = user.Name
	return u.userRepo.Update(existingUser)
}

func (u *userUsecase) Login(ctx context.Context, idToken string) (*domain.User, error) {
	// 1. Verify Token เหมือนใน Middleware
	token, err := u.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// 2. ดึง Email และเช็ค Verification
	email, _ := token.Claims["email"].(string)
	emailVerified, _ := token.Claims["email_verified"].(bool)

	if email == "" || !emailVerified {
		return nil, errors.New("email not verified or not found")
	}

	// 3. เช็คใน Database ว่า User นี้มีสิทธิ์เข้าใช้งานระบบเราไหม (Whitelist)
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		// กรณีไม่เจอใน DB แปลว่า Admin ยังไม่ได้เพิ่มอีเมลนี้เข้าระบบ
		return nil, errors.New("unauthorized: email not registered in system")
	}

	if user == nil {
		// นี่คือจุดที่ Email ไม่มีในตาราง users (Whitelist)
		return nil, errors.New("unauthorized: email not registered in system")
	}

	return user, nil
}

func (u *userUsecase) ListUsers() ([]*domain.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) GetUser(email string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *userUsecase) DeleteUser(email string) error {
	return u.userRepo.Delete(email)
}
