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
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Widget ID ไม่ถูกต้อง", nil)
	}

	var req SendCommandRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบข้อมูล JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลไม่ครบถ้วน",
			Errors:  fieldErrors,
		})
	}

	// 4. Business Logic
	widget, err := h.widgetUsecase.GetWidget(uint(widgetID))
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) { // หรือเพิ่ม ErrWidgetNotFound ใน domain
			return sendResponse(c, fiber.StatusNotFound, "ไม่พบ Widget ที่ระบุ", nil)
		}
		return sendResponse(c, fiber.StatusInternalServerError, "เกิดข้อผิดพลาดในการดึงข้อมูล Widget", nil)
	}

	correlationID := uuid.NewString()

	cmd := &domain.DeviceCommand{
		CapabilityType: widget.Capability.CapabilityType,
		ControlType:    widget.Capability.ControlType,
		Value:          req.Value,
		ReplyTopic:     fmt.Sprintf("devices/%s/command/response", correlationID),
	}

	if err := h.commandUsecase.SendCommand(cmd, req.Actor, uint(widgetID), correlationID); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถส่งคำสั่งไปยังอุปกรณ์ได้", nil)
	}

	// 5. Success Response
	return sendResponse(c, fiber.StatusAccepted, "ส่งคำสั่งเรียบร้อยแล้ว", fiber.Map{
		"correlation_id": correlationID,
	})
}