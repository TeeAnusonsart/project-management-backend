package http

import (
	"errors"
	"fmt"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
)

type WidgetHandler struct {
	usecase usecase.WidgetUsecase
}

func NewWidgetHandler(u usecase.WidgetUsecase) *WidgetHandler {
	return &WidgetHandler{usecase: u}
}

// --- Request DTOs for Validation ---

type CreateWidgetRequest struct {
	DeviceID     string `json:"device_id" validate:"required"`
	CapabilityID uint   `json:"capability_id" validate:"required"`
	WidgetStatus string `json:"widget_status" validate:"required,oneof=include exclude"`
}

type UpdateWidgetRequest struct {
	DeviceID     string `json:"device_id" validate:"required"`
	CapabilityID uint   `json:"capability_id" validate:"required"`
	WidgetStatus string `json:"widget_status" validate:"required,oneof=include exclude"`
}

type ChangeStatusRequest struct {
	WidgetStatus string `json:"widget_status" validate:"required,oneof=include exclude"`
}

type ChangeOrderRequest struct {
	WidgetOrders []uint `json:"widget_orders" validate:"required,min=0"`
}

// --- Handler Methods ---

// [GET] /widgets
func (h *WidgetHandler) ListWidgets(c *fiber.Ctx) error {
	status := c.Query("status")
	var widgets []*domain.Widget
	var err error

	if status != "" {
		widgets, err = h.usecase.GetWidgetByStatus(status)
	} else {
		widgets, err = h.usecase.ListWidgets()
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve widgets")
	}

	response := make([]WidgetResponse, 0)
	for _, w := range widgets {
		response = append(response, ToWidgetResponse(w))
	}

	return sendResponse(c, fiber.StatusOK, "Widgets retrieved successfully", response)
}

// [GET] /widgets/:widget_id
func (h *WidgetHandler) GetWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Widget ID format")
	}

	w, err := h.usecase.GetWidget(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) { // สามารถเปลี่ยนเป็น ErrWidgetNotFound ได้ตามโดเมนของคุณ
			return fiber.NewError(fiber.StatusNotFound, "Widget not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendResponse(c, fiber.StatusOK, "Widget retrieved successfully", ToWidgetResponse(w))
}

// [GET] /widgets/:widget_id/logs
func (h *WidgetHandler) GetLogs(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Widget ID format")
	}

	period := c.Query("period", "hour")

	logs, err := h.usecase.GetLogs(uint(id),period)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve logs")
	}

	response := make([]LogResponse, 0)
	for _, l := range logs {
		response = append(response, ToLogResponse(l))
	}

	return sendResponse(c, fiber.StatusOK, "Logs retrieved successfully", response)
}

// [POST] /widgets
func (h *WidgetHandler) CreateWidget(c *fiber.Ctx) error {
	var req CreateWidgetRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	widget := &domain.Widget{
		DeviceID:     req.DeviceID,
		CapabilityID: req.CapabilityID,
		WidgetStatus: req.WidgetStatus,
		Value:        "",
		WidgetOrder:  0,
	}

	if err := h.usecase.CreateWidget(widget); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create widget")
	}

	return sendResponse(c, fiber.StatusCreated, "Widget created successfully", nil)
}

// [PUT] /widgets/:widget_id
func (h *WidgetHandler) UpdateWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Widget ID format")
	}

	var req UpdateWidgetRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed",
			Errors:  fieldErrors,
		})
	}

	widget := &domain.Widget{
		ID:           uint(id),
		DeviceID:     req.DeviceID,
		CapabilityID: req.CapabilityID,
		WidgetStatus: req.WidgetStatus,
	}

	if err := h.usecase.UpdateWidget(widget); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update widget")
	}

	return sendResponse(c, fiber.StatusOK, "Widget updated successfully", nil)
}

// [PATCH] /widgets/:widget_id/status
func (h *WidgetHandler) ChangeStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Widget ID format")
	}

	var req ChangeStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed for status data",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.UpdateStatus(uint(id), req.WidgetStatus); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to change widget status")
	}

	return sendResponse(c, fiber.StatusOK, "Widget status updated successfully", nil)
}

// [PATCH] /rooms/:room_id/widgets/order
func (h *WidgetHandler) ChangeOrder(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	var req ChangeOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Validation failed for order data",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.ChangeOrder(uint(roomID), req.WidgetOrders); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to change widget order")
	}

	return sendResponse(c, fiber.StatusOK, "Widget order changed successfully", nil)
}

// [GET] /rooms/:room_id/widgets
func (h *WidgetHandler) ListByRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Room ID format")
	}

	status := c.Query("status")
	widgets, err := h.usecase.ListWidgetsByRoomWithStatus(uint(roomID), status)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve widgets for the room")
	}

	response := make([]WidgetResponse, 0)
	for _, w := range widgets {
		fmt.Print(w.ID)
		response = append(response, ToWidgetResponse(w))
	}

	return sendResponse(c, fiber.StatusOK, "Widgets for the room retrieved successfully", response)
}

// [DELETE] /widgets/:widget_id
func (h *WidgetHandler) DeleteWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Widget ID format")
	}

	if err := h.usecase.DeleteWidget(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete widget")
	}

	return sendResponse(c, fiber.StatusOK, "Widget deleted successfully", nil)
}

// --- Response Helpers & DTOs ---

type LogResponse struct {
	Value     string    `json:"value"`
	EventType string    `json:"event_type"`
	Actor     string    `json:"actor"`
	CreatedAt time.Time `json:"created_at"`
}

type WidgetResponse struct {
	WidgetID     uint           `json:"widget_id"`
	WidgetOrder  uint           `json:"widget_order"`
	WidgetStatus string         `json:"widget_status"`
	Value        string         `json:"value"`
	Device       *DeviceDTO     `json:"device,omitempty"`
	Capability   *CapabilityDTO `json:"capability,omitempty"`
}

type DeviceDTO struct {
	DeviceID            string    `json:"device_id"`
	DeviceName          string    `json:"device_name"`
	DeviceLastHeartbeat time.Time `json:"device_last_heartbeat"`
	DeviceType          string    `json:"device_type"`
}

type CapabilityDTO struct {
	CapabilityID   uint   `json:"capability_id"`
	CapabilityType string `json:"capability_type"`
	ControlType    string `json:"control_type"`
}

func ToLogResponse(l *domain.Log) LogResponse {
	return LogResponse{
		Value:     l.Value,
		EventType: l.EventType,
		Actor:     l.Actor,
		CreatedAt: l.CreatedAt,
	}
}

func ToWidgetResponse(w *domain.Widget) WidgetResponse {
	response := WidgetResponse{
		WidgetID:     w.ID,
		WidgetOrder:  w.WidgetOrder,
		WidgetStatus: w.WidgetStatus,
		Value:        w.Value,
	}

	if w.Device != nil {
		response.Device = &DeviceDTO{
			DeviceID:            w.Device.DeviceID,
			DeviceLastHeartbeat: w.Device.LastHeartbeat,
			DeviceName:          w.Device.DeviceName,
			DeviceType:          w.Device.DeviceType,
		}
	}

	if w.Capability != nil {
		response.Capability = &CapabilityDTO{
			CapabilityID:   w.Capability.ID,
			CapabilityType: w.Capability.CapabilityType,
			ControlType:    w.Capability.ControlType,
		}
	}

	return response
}