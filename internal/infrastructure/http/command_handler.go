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
}

func NewCommandHandler(uc usecase.CommandUsecase) *CommandHandler {
	return &CommandHandler{commandUsecase: uc}
}

func (h *CommandHandler) SendCommand(c *fiber.Ctx) error {
	widgetIDParam := c.Params("widgetId")
	widgetID, err := strconv.ParseUint(widgetIDParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid widget id")
	}

	var req struct {
		CapabilityID string `json:"capability_id"`
		Value   uint   `json:"value"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	correlationID := uuid.NewString()

	cmd := &domain.DeviceCommand{
		WidgetID:      uint(widgetID),
		CapabilityID:       req.CapabilityID,
		Value:         req.Value,
		CorrelationID: correlationID,
		ReplyTopic:    "devices/reply/" + correlationID,
	}

	if err := h.commandUsecase.SendCommand(cmd); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "success",
	})
}
