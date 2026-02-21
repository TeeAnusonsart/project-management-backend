package http

import (
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

func NewDeviceHandler(u usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: u}
}

func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	// status := c.Query("connected")
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
	// if status == "ture" {
	// 	fmt.Println(status)
	// 	devices, err = h.usecase.GetPairedDevice()
	// } else if status == "false" {
	// 	devices, err = h.usecase.GetUnpairDevice()
	// } else {
	// 	devices, err = h.usecase.ListDevices()
	// }
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []DeviceResponse
	for _, d := range devices {
		response = append(response, ToDeviceResponse(d))
	}

	return c.JSON(fiber.Map{"data": response})
}

func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")

	device, err := h.usecase.GetDevice(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(ToDeviceResponse(device))
}

func (h *DeviceHandler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")

	var req UpdateDeviceRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.UpdateDevice(id, req.DeviceName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "แก้ไขข้อมูลของอุปกรณ์เรียบร้อยแล้ว",
	})
}

func (h *DeviceHandler) PairDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")

	var req PairDeviceRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.PairDevice(id, req.DeviceKey); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เชื่อมต่ออุปกรณ์เรียบร้อยแล้ว",
	})
}

func (h *DeviceHandler) UnpairDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")

	if err := h.usecase.UnpairDevice(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "ยกเลิกการเชื่อมต่ออุปกรณ์เรียบร้อยแล้ว",
	})
}

type UpdateDeviceRequest struct {
	DeviceName string `json:"device_name"`
}

type PairDeviceRequest struct {
	DeviceKey string `json:"device_key"`
}

type DeviceResponse struct {
	DeviceID            string    `json:"device_id"`
	DeviceLastHeartbeat time.Time `json:"device_last_heartbeat"`
	DeviceName          string    `json:"device_name"`
	DeviceType          string    `json:"device_type"`
}

func ToDeviceResponse(d *domain.DeviceSummary) DeviceResponse {
	return DeviceResponse{
		DeviceID:            d.DeviceID,
		DeviceLastHeartbeat: d.LastHeartbeat,
		DeviceName:          d.DeviceName,
		DeviceType:          d.DeviceType,
	}
}
