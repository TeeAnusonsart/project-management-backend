package http

import (
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

func NewCommandHandler(uc usecase.CommandUsecase,wc usecase.WidgetUsecase) *CommandHandler {
	return &CommandHandler{commandUsecase: uc,widgetUsecase: wc}
}

func (h *CommandHandler) SendCommand(c *fiber.Ctx) error {
	widgetIDParam := c.Params("widgetId")
	widgetID, err := strconv.ParseUint(widgetIDParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid widget id")
	}

	var req struct {
		CapabilityID uint `json:"capability_id"`
		Value        uint   `json:"value"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	widget, err := h.widgetUsecase.GetWidget(uint(widgetID))

	correlationID := uuid.NewString()

	cmd := &domain.DeviceCommand{
		CapabilityType: widget.Capability.CapabilityType,
		ControlType:    widget.Capability.ControlType,
		Value:          req.Value,
		ReplyTopic:     "devices/reply/" + correlationID,
	}

	if err := h.commandUsecase.SendCommand(cmd, uint(widgetID), correlationID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "success",
	})
}
