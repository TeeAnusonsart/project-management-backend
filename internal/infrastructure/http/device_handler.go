package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"project-home-iot/internal/core/usecase"
)

type DeviceHandler struct {
	usecase usecase.DeviceUsecase
}

func NewDeviceHandler(u usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{usecase: u}
}

func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	devices, err := h.usecase.ListDevices()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"data": devices})
}

func (h *DeviceHandler) ListDevices(c *fiber.Ctx) error {
	devices, err := h.usecase.ListDevices()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"data": devices})
}

func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("device_id"), 10, 64)

	device, err := h.usecase.GetDevice(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(device)
}

func (h *DeviceHandler) UpdateDevice(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("device_id"), 10, 64)

	var req struct {
		DeviceName string `json:"device_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.UpdateDevice(uint(id), req.DeviceName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "แก้ไขข้อมูลของอุปกรณ์เรียบร้อยแล้ว",
	})
}

func (h *DeviceHandler) PairDevice(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("device_id"), 10, 64)

	var req struct {
		DeviceKey string `json:"device_key"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.PairDevice(uint(id), req.DeviceKey); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เชื่อมต่ออุปกรณ์เรียบร้อยแล้ว",
	})
}

func (h *DeviceHandler) UnpairDevice(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("device_id"), 10, 64)

	if err := h.usecase.UnpairDevice(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "ยกเลิกการเชื่อมต่ออุปกรณ์เรียบร้อยแล้ว",
	})
}
