package dtos

type CreateRoomRequest struct {
	RoomName string `json:"room_name"`
}

type UpdateRoomRequest struct {
	RoomName string `json:"room_name"`
}

type AddDeviceToRoomRequest struct {
	DeviceID uint `json:"device_id"`
}
