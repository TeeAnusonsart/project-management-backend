package http

import (
	"errors"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"time"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

// --- Handler ---

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

func NewDeviceHandler(u usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: u}
}

// [GET] /devices (รองรับ Query ?connected=true/false)
func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	var devices []*domain.DeviceSummary
	var err error
	connected := c.Query("connected")

	if connected == "" {
		devices, err = h.usecase.ListDevices()
	} else {
		isConnected, parseErr := strconv.ParseBool(connected)
		if parseErr != nil {
			return sendResponse(
				c,
				fiber.StatusBadRequest,
				"connected must be true or false",
				nil,
			)
		}

		if isConnected {
			devices, err = h.usecase.GetPairedDevice()
		} else {
			devices, err = h.usecase.GetUnpairDevice()
		}
	}
	
	

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := make([]DeviceResponse, 0)
	for _, d := range devices {
		response = append(response, ToDeviceResponse(d))
	}

	return sendResponse(c, fiber.StatusOK, "Devices retrieved successfully", response)
}

// [GET] /devices/:device_id
func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Device ID is required")
	}

	device, err := h.usecase.GetDevice(id)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Device not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get device")
	}

	return sendResponse(c, fiber.StatusOK, "Device retrieved successfully", ToDeviceResponse(device))
}

// [PUT] /devices/:device_id
func (h *DeviceHandler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Device ID is required")
	}

	var req UpdateDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// ใช้ Validation แบบส่ง fieldErrors กลับไป (เหมือน RoomHandler)
	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.UpdateDevice(id, req.DeviceName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendResponse(c, fiber.StatusOK, "Device updated successfully", nil)
}

// [POST] /devices/:device_id/pair
func (h *DeviceHandler) PairDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Device ID is required")
	}

	var req PairDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// ใช้ Validation แบบส่ง fieldErrors กลับไป (เหมือน RoomHandler)
	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	_, err := h.usecase.GetDevice(id)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Device not found")
		}
	}

	err = h.usecase.PairDevice(id, req.DeviceKey)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDeviceAlreadyPaired):
			return fiber.NewError(fiber.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrDeviceNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrInvalidDeviceKey):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case errors.Is(err, domain.ErrPairFailed):
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	return sendResponse(c, fiber.StatusOK, "Device paired successfully", nil)
}

// [POST] /devices/:device_id/unpair
func (h *DeviceHandler) UnpairDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Device ID is required")
	}

	if err := h.usecase.UnpairDevice(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return sendResponse(c, fiber.StatusOK, "Device unpaired successfully", nil)
}

// --- Request / Response DTOs ---

type UpdateDeviceRequest struct {
	DeviceName string `json:"device_name" validate:"required,min=2,max=50"`
}

type PairDeviceRequest struct {
	DeviceKey string `json:"device_key" validate:"required,min=4"`
}

type DeviceResponse struct {
	DeviceID            string     `json:"device_id"`
	DeviceName          string     `json:"device_name"`
	DeviceType          string     `json:"device_type"`
	DeviceLastHeartbeat *time.Time `json:"device_last_heartbeat"`
}

func ToDeviceResponse(d *domain.DeviceSummary) DeviceResponse {
	return DeviceResponse{
		DeviceID:            d.DeviceID,
		DeviceName:          d.DeviceName,
		DeviceType:          d.DeviceType,
		DeviceLastHeartbeat: &d.LastHeartbeat,
	}
}