package http

import (
	"errors"
	"fmt"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)



// --- Handler ---

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

var validate = validator.New()

func NewDeviceHandler(u usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: u}
}

// ฟังก์ชันช่วยจัดการ Validation Error ให้เป็นข้อความที่อ่านง่าย
func validateStruct(s interface{}) string {
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return fmt.Sprintf("ฟิลด์ %s จำเป็นต้องระบุ", err.Field())
			case "min":
				return fmt.Sprintf("ฟิลด์ %s ต้องมีความยาวอย่างน้อย %s ตัวอักษร", err.Field(), err.Param())
			case "max":
				return fmt.Sprintf("ฟิลด์ %s ต้องมีความยาวไม่เกิน %s ตัวอักษร", err.Field(), err.Param())
			}
		}
		return err.Error()
	}
	return ""
}

func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	var devices []*domain.DeviceSummary
	var err error
	connected := c.Query("connected")

	if connected == "" {
		devices, err = h.usecase.ListDevices()
	} else {
		isConnected := c.QueryBool("connected")
		if isConnected {
			devices, err = h.usecase.GetPairedDevice()
		} else {
			devices, err = h.usecase.GetUnpairDevice()
		}
	}

	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	response := make([]DeviceResponse, 0)
	for _, d := range devices {
		response = append(response, ToDeviceResponse(d))
	}

	return sendResponse(c, fiber.StatusOK, "Devices retrieved successfully", response)
}

func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	if id == "" {
		return sendResponse(c, fiber.StatusBadRequest, "Device ID is required", nil)
	}

	device, err := h.usecase.GetDevice(id)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			return sendResponse(c, fiber.StatusNotFound, "Device not found", nil)
		}
		return sendResponse(c, fiber.StatusInternalServerError, "Failed to get device", nil)
	}

	return sendResponse(c, fiber.StatusOK, "Device retrieved successfully", ToDeviceResponse(device))
}

func (h *DeviceHandler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	var req UpdateDeviceRequest

	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validation
	if errMsg := validateStruct(req); errMsg != "" {
		return sendResponse(c, fiber.StatusUnprocessableEntity, errMsg, nil)
	}

	if err := h.usecase.UpdateDevice(id, req.DeviceName); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return sendResponse(c, fiber.StatusOK, "แก้ไขข้อมูลของอุปกรณ์เรียบร้อยแล้ว", nil)
}

func (h *DeviceHandler) PairDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	var req PairDeviceRequest

	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	// Validation
	if errMsg := validateStruct(req); errMsg != "" {
		return sendResponse(c, fiber.StatusUnprocessableEntity, errMsg, nil)
	}

	err := h.usecase.PairDevice(id, req.DeviceKey)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDeviceAlreadyPaired):
			return sendResponse(c, fiber.StatusConflict, err.Error(), nil)
		case errors.Is(err, domain.ErrDeviceNotFound):
			return sendResponse(c, fiber.StatusNotFound, err.Error(), nil)
		default:
			return sendResponse(c, fiber.StatusInternalServerError, "Internal server error", nil)
		}
	}

	return sendResponse(c, fiber.StatusOK, "เชื่อมต่ออุปกรณ์เรียบร้อยแล้ว", nil)
}

func (h *DeviceHandler) UnpairDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")
	if err := h.usecase.UnpairDevice(id); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return sendResponse(c, fiber.StatusOK, "ยกเลิกการเชื่อมต่ออุปกรณ์เรียบร้อยแล้ว", nil)
}



type UpdateDeviceRequest struct {
	DeviceName string `json:"device_name" validate:"required,min=2,max=50"`
}

type PairDeviceRequest struct {
	DeviceKey string `json:"device_key" validate:"required,min=4"`
}

type DeviceResponse struct {
	DeviceID            string    `json:"id"` // ปรับให้ตรงตามตัวอย่างที่คุณต้องการ
	DeviceName          string    `json:"name"`
	DeviceType          string    `json:"type"`
	DeviceLastHeartbeat *time.Time `json:"lastHeartbeatAt"`
}

func ToDeviceResponse(d *domain.DeviceSummary) DeviceResponse {
	return DeviceResponse{
		DeviceID:            d.DeviceID,
		DeviceName:          d.DeviceName,
		DeviceType:          d.DeviceType,
		DeviceLastHeartbeat: &d.LastHeartbeat,
	}
}