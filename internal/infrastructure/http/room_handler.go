package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"project-home-iot/internal/core/usecase"
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
	return c.JSON(fiber.Map{"data": rooms})
}

func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("room_id"), 10, 64)

	room, err := h.usecase.GetRoom(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(fiber.Map{
		"room_id":   room.ID,
		"room_name": room.Name,
	})
}

func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req struct {
		RoomName string `json:"room_name"`
	}

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

	var req struct {
		RoomName string `json:"room_name"`
	}

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

	var req struct {
		DeviceID uint `json:"device_id"`
	}

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

	return c.JSON(fiber.Map{"data": devices})
}
