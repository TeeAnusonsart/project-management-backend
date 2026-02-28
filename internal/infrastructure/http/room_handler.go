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
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูลห้องได้", nil)
	}

	response := make([]RoomResponse, 0)
	for _, r := range rooms {
		response = append(response, ToRoomResponse(r))
	}

	return sendResponse(c, fiber.StatusOK, "เรียกดูข้อมูลสำเร็จ", response)
}

// [GET] /rooms/:room_id
func (h *RoomHandler) GetRoom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	room, err := h.usecase.GetRoom(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) { // ปรับเป็น ErrRoomNotFound หากมีใน domain
			return sendResponse(c, fiber.StatusNotFound, "ไม่พบข้อมูลห้อง", nil)
		}
		return sendResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return sendResponse(c, fiber.StatusOK, "เรียกดูข้อมูลสำเร็จ", ToRoomResponse(room))
}

// [POST] /rooms
func (h *RoomHandler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลไม่ถูกต้อง",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.CreateRoom(req.RoomName); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถสร้างห้องได้", nil)
	}

	return sendResponse(c, fiber.StatusCreated, "เพิ่มห้องใหม่เรียบร้อยแล้ว", nil)
}

// [PUT] /rooms/:room_id
func (h *RoomHandler) UpdateRoom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	var req UpdateRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลไม่ถูกต้อง",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.UpdateRoom(uint(id), req.RoomName); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถแก้ไขข้อมูลห้องได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "แก้ไขข้อมูลของห้องเรียบร้อยแล้ว", nil)
}

// [DELETE] /rooms/:room_id
func (h *RoomHandler) DeleteRoom(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	if err := h.usecase.DeleteRoom(uint(id)); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถลบห้องได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "นำห้องออกเรียบร้อยแล้ว", nil)
}

// [POST] /rooms/:room_id/devices
func (h *RoomHandler) AddDevice(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	var req AddDeviceToRoomRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลอุปกรณ์ไม่ถูกต้อง",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.AddDeviceToRoom(uint(roomID), req.DeviceID); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถเพิ่มอุปกรณ์เข้าห้องได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "เพิ่มอุปกรณ์ภายในห้องเรียบร้อยแล้ว", nil)
}

// [GET] /rooms/:room_id/devices
func (h *RoomHandler) ListDevices(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	devices, err := h.usecase.ListDevicesInRoom(uint(roomID))
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูลอุปกรณ์ในห้องได้", nil)
	}

	response := make([]DeviceResponse, 0)
	for _, d := range devices {
		response = append(response, ToDeviceResponse(d))
	}

	return sendResponse(c, fiber.StatusOK, "เรียกดูอุปกรณ์ภายในห้องสำเร็จ", response)
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