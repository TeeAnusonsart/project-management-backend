package mappers

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/http/dtos"
)

func ToRoomResponse(r *domain.Room) dtos.RoomResponse {
	return dtos.RoomResponse{
		RoomID:   r.ID,
		RoomName: r.Name,
	}
}
