package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/http/dtos"
)

func ToUserResponse(u *domain.User) dtos.UserResponse {
	return dtos.UserResponse{
		UserID:      u.ID,
		Username:    u.Username,
		Name:        u.Name,
		Email:       u.Email,
		ProfilePath: u.ProfilePath,
		Role:        string(u.Role),
	}
}
