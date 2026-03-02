package http

import (
	"errors"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	usecase usecase.RoomUsecase
}

func NewRoomHandler(u usecase.RoomUsecase) *RoomHandler {
	return &RoomHandler{usecase: u}
}

// --- Request DTOs ---

type CreateRoomRequest struct {
	RoomName string `json:"room_name" validate:"required,min=2,max=100"`
}

type UpdateRoomRequest struct {
	RoomName string `json:"room_name" validate:"required,min=2,max=100"`
}

type AddDeviceToRoomRequest struct {
	DeviceID string `json:"device_id" validate:"required"`
}

// --- Handler Methods ---

// [GET] /rooms
func (h *RoomHandler) ListRooms(c *fiber.Ctx) error {
	rooms, err := h.usecase.ListRooms()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve rooms")
	}

	response := make([]RoomResponse, 0)
	for _, r := range rooms {
		response = append(response, ToRoomResponse(r))
	}

	return sendResponse(c, fiber.StatusOK, "Rooms retrieved successfully", response)
}

// [GET] /rooms/:room_id
func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	room, err := h.usecase.GetRoom(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoomNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Room not found")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve room")
		}
	}

	return sendResponse(c, fiber.StatusOK, "Room retrieved successfully", ToRoomResponse(room))
}

// [POST] /rooms
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	err := h.usecase.CreateRoom(req.RoomName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoomExist):
			return fiber.NewError(fiber.StatusConflict, err.Error())

		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create room")
		}
	}

	return sendResponse(c, fiber.StatusCreated, "Room created successfully", nil)
}

// [PUT] /rooms/:room_id
func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	var req UpdateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.UpdateRoom(uint(id), req.RoomName); err != nil {
		switch {
		case errors.Is(err, domain.ErrRoomNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Room not found")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update room")
		}
	}

	return sendResponse(c, fiber.StatusOK, "Room updated successfully", nil)
}

// [DELETE] /rooms/:room_id
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	if err := h.usecase.DeleteRoom(uint(id)); err != nil {
		switch {
		case errors.Is(err, domain.ErrRoomNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete room")
		}
	}

	return sendResponse(c, fiber.StatusOK, "Room deleted successfully", nil)
}

// [POST] /rooms/:room_id/devices
func (h *RoomHandler) AddDevice(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	var req AddDeviceToRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed for device data",
			Errors:  fieldErrors,
		})
	}

	err = h.usecase.AddDeviceToRoom(uint(roomID), req.DeviceID)
	if err != nil {
		switch err {
		case domain.ErrRoomNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case domain.ErrDeviceNotFound:
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to add device to room")
		}
	}

	return sendResponse(c, fiber.StatusOK, "Device added to room successfully", nil)
}

// [GET] /rooms/:room_id/devices
func (h *RoomHandler) ListDevices(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	devices, err := h.usecase.ListDevicesInRoom(uint(roomID))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRoomNotFound):
			return fiber.NewError(fiber.StatusNotFound, "Room not found")
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve devices in room")
		}
	}

	response := make([]DeviceResponse, 0)
	for _, d := range devices {
		response = append(response, ToDeviceResponse(d))
	}

	return sendResponse(c, fiber.StatusOK, "Devices in room retrieved successfully", response)
}

// --- Response Helpers & DTOs ---

type RoomResponse struct {
	RoomID   uint             `json:"room_id"`
	RoomName string           `json:"room_name"`
	Devices  []DeviceResponse `json:"devices,omitempty"`
}

func ToRoomResponse(r *domain.Room) RoomResponse {
	devices := make([]DeviceResponse, 0)
	for _, d := range r.Devices {
		devices = append(devices, ToDeviceResponse(&d))
	}

	return RoomResponse{
		RoomID:   r.ID,
		RoomName: r.Name,
		Devices:  devices,
	}
}
