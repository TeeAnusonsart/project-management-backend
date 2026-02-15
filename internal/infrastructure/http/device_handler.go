package http

import (
	"fmt"
	"project-home-iot/internal/core/usecase"
	"project-home-iot/internal/infrastructure/http/dtos"
	httpmapper "project-home-iot/internal/infrastructure/http/mappers"

	"github.com/gofiber/fiber/v2"
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

	var response []dtos.DeviceResponse
	for _, d := range devices {
		response = append(response, httpmapper.ToDeviceResponse(d))
	}

	return c.JSON(fiber.Map{"data": response})
}


func (h *DeviceHandler) GetDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")

	device, err := h.usecase.GetDevice(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(httpmapper.ToDeviceResponse(device))
}

func (h *DeviceHandler) UpdateDevice(c *fiber.Ctx) error {
	id := c.Params("device_id")

	var req dtos.UpdateDeviceRequest

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

	fmt.Printf("Pairing device with ID: %s\n", id)

	var req dtos.PairDeviceRequest

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
