package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"project-home-iot/internal/core/usecase"
	"project-home-iot/internal/core/domain"
)

type RoomHandler struct {
	usecase usecase.RoomUsecase
}

func NewRoomHandler(u usecase.RoomUsecase) *RoomHandler {
	return &RoomHandler{usecase: u}
}

func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	rooms, err := h.usecase.ListRooms()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []RoomResponse
	for _, r := range rooms {
		response = append(response, ToRoomResponse(r))
	}

	return c.JSON(fiber.Map{"data": response})
}

func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("room_id"), 10, 64)

	room, err := h.usecase.GetRoom(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(ToRoomResponse(room))
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.CreateRoom(req.RoomName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เพิ่มห้องใหม่เรียบร้อยแล้ว",
	})
}

func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("room_id"), 10, 64)
	var req UpdateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.UpdateRoom(uint(id), req.RoomName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "แก้ไขข้อมูลของห้องเรียบร้อยแล้ว",
	})
}

func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("room_id"), 10, 64)

	if err := h.usecase.DeleteRoom(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "นำห้องออกเรียบร้อยแล้ว",
	})
}

func (h *RoomHandler) AddDevice(c *fiber.Ctx) error {
	roomID, _ := strconv.ParseUint(c.Params("room_id"), 10, 64)

	var req AddDeviceToRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.AddDeviceToRoom(uint(roomID), req.DeviceID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เพิ่มอุปกรณ์ภายในห้องเรียบร้อยแล้ว",
	})
}

func (h *RoomHandler) ListDevices(c *fiber.Ctx) error {
	roomID, _ := strconv.ParseUint(c.Params("room_id"), 10, 64)

	devices, err := h.usecase.ListDevicesInRoom(uint(roomID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []DeviceResponse
	for _, d := range devices {
		response = append(response, ToDeviceResponse(d))
	}

	return c.JSON(fiber.Map{"data": response})
}

func ToRoomResponse(r *domain.Room) RoomResponse {
	return RoomResponse{
		RoomID:   r.ID,
		RoomName: r.Name,
	}
}


type CreateRoomRequest struct {
    RoomName string `json:"room_name"`
}

type UpdateRoomRequest struct {
    RoomName string `json:"room_name"`
}

type AddDeviceToRoomRequest struct {
    DeviceID uint `json:"device_id"`
}

type RoomResponse struct {
    RoomID   uint   `json:"room_id"`
    RoomName string `json:"room_name"`
}