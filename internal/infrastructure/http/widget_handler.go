package http

import (
	"errors"
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
	WidgetStatus string `json:"widget_status" validate:"required,oneof=active inactive"`
}

type UpdateWidgetRequest struct {
	DeviceID     string `json:"device_id" validate:"required"`
	CapabilityID uint   `json:"capability_id" validate:"required"`
	WidgetStatus string `json:"widget_status" validate:"required,oneof=active inactive"`
}

type ChangeStatusRequest struct {
	WidgetStatus string `json:"widget_status" validate:"required,oneof=active inactive"`
}

type ChangeOrderRequest struct {
	WidgetOrders []uint `json:"widget_orders" validate:"required,min=1"`
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
		return sendResponse(c, fiber.StatusInternalServerError, "เกิดข้อผิดพลาดในการดึงข้อมูล", nil)
	}

	response := make([]WidgetResponse, 0)
	for _, w := range widgets {
		response = append(response, ToWidgetResponse(w))
	}

	return sendResponse(c, fiber.StatusOK, "เรียกดูข้อมูลสำเร็จ", response)
}

// [GET] /widgets/:widget_id
func (h *WidgetHandler) GetWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Widget ID ไม่ถูกต้อง", nil)
	}

	w, err := h.usecase.GetWidget(uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) { // สามารถเปลี่ยนเป็น ErrWidgetNotFound ได้
			return sendResponse(c, fiber.StatusNotFound, "ไม่พบข้อมูล Widget", nil)
		}
		return sendResponse(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return sendResponse(c, fiber.StatusOK, "เรียกดูข้อมูลสำเร็จ", ToWidgetResponse(w))
}

// [GET] /widgets/:widget_id/logs
func (h *WidgetHandler) GetLogs(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Widget ID ไม่ถูกต้อง", nil)
	}

	logs, err := h.usecase.GetLogs(uint(id))
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูล Logs ได้", nil)
	}

	response := make([]LogResponse, 0)
	for _, l := range logs {
		response = append(response, ToLogResponse(l))
	}

	return sendResponse(c, fiber.StatusOK, "เรียกดู Logs สำเร็จ", response)
}

// [POST] /widgets
func (h *WidgetHandler) CreateWidget(c *fiber.Ctx) error {
	var req CreateWidgetRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลไม่ถูกต้อง",
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
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถสร้าง Widget ได้", nil)
	}

	return sendResponse(c, fiber.StatusCreated, "เพิ่ม widget ใหม่เรียบร้อยแล้ว", nil)
}

// [PUT] /widgets/:widget_id
func (h *WidgetHandler) UpdateWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Widget ID ไม่ถูกต้อง", nil)
	}

	var req UpdateWidgetRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลไม่ถูกต้อง",
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
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถแก้ไขข้อมูลได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "แก้ไขข้อมูลของ widget เรียบร้อยแล้ว", nil)
}

// [PATCH] /widgets/:widget_id/status
func (h *WidgetHandler) ChangeStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Widget ID ไม่ถูกต้อง", nil)
	}

	var req ChangeStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลสถานะไม่ถูกต้อง",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.UpdateStatus(uint(id), req.WidgetStatus); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถเปลี่ยนสถานะได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "การดำเนินการเสร็จสิ้น", nil)
}

// [PATCH] /rooms/:room_id/widgets/order
func (h *WidgetHandler) ChangeOrder(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	var req ChangeOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "JSON ไม่ถูกต้อง", nil)
	}

	if fieldErrors := validateStruct(req); len(fieldErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(domain.ErrorResponse{
			Status:  "error",
			Code:    fiber.StatusUnprocessableEntity,
			Message: "ข้อมูลลำดับไม่ถูกต้อง",
			Errors:  fieldErrors,
		})
	}

	if err := h.usecase.ChangeOrder(uint(roomID), req.WidgetOrders); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถเปลี่ยนลำดับได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "เปลี่ยนลำดับ widget เรียบร้อยแล้ว", nil)
}

// [GET] /rooms/:room_id/widgets
func (h *WidgetHandler) ListByRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Room ID ไม่ถูกต้อง", nil)
	}

	status := c.Query("status")
	widgets, err := h.usecase.ListWidgetsByRoomWithStatus(uint(roomID), status)
	if err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถดึงข้อมูลได้", nil)
	}

	response := make([]WidgetResponse, 0)
	for _, w := range widgets {
		response = append(response, ToWidgetResponse(w))
	}

	return sendResponse(c, fiber.StatusOK, "ดึงข้อมูล Widget ตามห้องสำเร็จ", response)
}

// [DELETE] /widgets/:widget_id
func (h *WidgetHandler) DeleteWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return sendResponse(c, fiber.StatusBadRequest, "รูปแบบ Widget ID ไม่ถูกต้อง", nil)
	}

	if err := h.usecase.DeleteWidget(uint(id)); err != nil {
		return sendResponse(c, fiber.StatusInternalServerError, "ไม่สามารถลบ Widget ได้", nil)
	}

	return sendResponse(c, fiber.StatusOK, "นำ widget ออกเรียบร้อยแล้ว", nil)
}

// --- Response Helpers & DTOs ---

type LogResponse struct {
	Value     string    `json:"value"`
	EventType string    `json:"event_type"`
	Actor     string    `json:"actor"`
	CreatedAt time.Time `json:"created_at"`
}

type WidgetResponse struct {
	WidgetID     uint          `json:"widget_id"`
	WidgetOrder  uint          `json:"widget_order"`
	WidgetStatus string        `json:"widget_status"`
	Value        string        `json:"value"`
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