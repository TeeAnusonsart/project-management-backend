package http

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
)

type CommandHandler struct {
	commandUsecase usecase.CommandUsecase
	widgetUsecase  usecase.WidgetUsecase
}

func NewCommandHandler(uc usecase.CommandUsecase, wc usecase.WidgetUsecase) *CommandHandler {
	return &CommandHandler{commandUsecase: uc, widgetUsecase: wc}
}

type SendCommandRequest struct {
	Actor string `json:"actor" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (h *CommandHandler) SendCommand(c *fiber.Ctx) error {
	widgetIDParam := c.Params("widgetId")
	widgetID, err := strconv.ParseUint(widgetIDParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Widget ID format")
	}

	var req SendCommandRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	// ใช้ validateStruct จากไฟล์ validator.go ที่เราสร้างไว้ร่วมกัน
	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลไม่ครบถ้วน",
			Errors:  fieldErrors,
		})
	}

	widget, err := h.widgetUsecase.GetWidget(uint(widgetID))
	if err != nil {
		// เช็คว่าเป็น Error แบบ "หาไม่เจอ" หรือไม่ (อย่าลืมประกาศ ErrWidgetNotFound ใน package domain ด้วยนะครับ)
		if errors.Is(err, domain.ErrWidgetNotFound) { 
			return fiber.NewError(fiber.StatusNotFound, "Widget not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve widget data")
	}

	

	correlationID := uuid.NewString()

	cmd := &domain.DeviceCommand{
		CapabilityType: widget.Capability.CapabilityType,
		ControlType:    widget.Capability.ControlType,
		Value:          req.Value,
		ReplyTopic:     fmt.Sprintf("devices/%s/command/response", correlationID),
	}

	if err := h.commandUsecase.SendCommand(cmd, req.Actor, uint(widgetID), correlationID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send command to the device")
	}

	// Success Response
	return sendResponse(c, fiber.StatusAccepted, "Command sent successfully", fiber.Map{
		"correlation_id": correlationID,
	})
}