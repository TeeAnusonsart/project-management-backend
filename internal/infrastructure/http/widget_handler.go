package http

import (
	"strconv"

	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"time"
	"github.com/gofiber/fiber/v2"
)

type WidgetHandler struct {
	usecase usecase.WidgetUsecase
}

func NewWidgetHandler(u usecase.WidgetUsecase) *WidgetHandler {
	return &WidgetHandler{usecase: u}
}

func (h *WidgetHandler) ListWidgets(c *fiber.Ctx) error {
	status := c.Query("status")
	var widgets []*domain.Widget
	var err error

	if status != "" {
        widgets, err = h.usecase.GetWidgetByStatus(status)
    } else {
        widgets, err = h.usecase.ListWidgets()
    }

	// widgets, err := h.usecase.ListWidgets()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []WidgetResponse
	for _, w := range widgets {
		response = append(response, ToWidgetResponse(w))
	}

	return c.JSON(fiber.Map{"data": response})
}

func (h *WidgetHandler) GetWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid widget id")
	}

	w, err := h.usecase.GetWidget(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(ToWidgetResponse(w))
}


func (h *WidgetHandler) CreateWidget(c *fiber.Ctx) error {
	var req struct {
		DeviceID     string   `json:"device_id"`
		CapabilityID uint   `json:"capability_id"`
		WidgetStatus string `json:"widget_status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	widget := &domain.Widget{
		DeviceID:     req.DeviceID,
		CapabilityID: req.CapabilityID,
		WidgetStatus: req.WidgetStatus,
		Value:        0,
		WidgetOrder:  0,
	}

	if err := h.usecase.CreateWidget(widget); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เพิ่ม widget ใหม่เรียบร้อยแล้ว",
	})
}


func (h *WidgetHandler) UpdateWidget(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("widget_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid widget id")
	}

	var req struct {
		DeviceID     string   `json:"device_id"`
		CapabilityID uint   `json:"capability_id"`
		WidgetStatus string `json:"widget_status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	widget := &domain.Widget{
		ID:           uint(id),
		DeviceID:     req.DeviceID,
		CapabilityID: req.CapabilityID,
		WidgetStatus: req.WidgetStatus,
	}

	if err := h.usecase.UpdateWidget(widget); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "แก้ไขข้อมูลของ widget เรียบร้อยแล้ว",
	})
}

func (h *WidgetHandler) ChangeStatus(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("widget_id"), 10, 64)

	var req struct {
		WidgetStatus string `json:"widget_status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.UpdateStatus(uint(id), req.WidgetStatus); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "การดำเนินการเสร็จสิ้น",
	})
}

func (h *WidgetHandler) ChangeOrder(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid room id")
	}

	var req struct {
		WidgetOrders []uint `json:"widget_orders"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.usecase.ChangeOrder(uint(roomID), req.WidgetOrders); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "เปลี่ยนลำดับ widget เรียบร้อยแล้ว",
	})
}

func (h *WidgetHandler) ListByRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("room_id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid room id")
	}

	widgets, err := h.usecase.ListWidgetsByRoom(uint(roomID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []WidgetResponse
	for _, w := range widgets {
		response = append(response, ToWidgetResponse(w))
	}

	return c.JSON(fiber.Map{"data": response})
}

func (h *WidgetHandler) DeleteWidget(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("widget_id"), 10, 64)

	if err := h.usecase.DeleteWidget(uint(id)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "นำ widget ออกเรียบร้อยแล้ว",
	})
}

func ToWidgetResponse(w *domain.Widget) WidgetResponse {

	response := WidgetResponse{
		WidgetID:     w.ID,
		WidgetOrder:  w.WidgetOrder,
		WidgetStatus: w.WidgetStatus,
		Value:        w.Value,
	}

	if w.Device != nil {
		response.Device = DeviceDTO{
			DeviceID:   w.Device.DeviceID,
			DeviceLastHeartbeat: w.Device.LastHeartbeat,
			DeviceName: w.Device.DeviceName,
			DeviceType: w.Device.DeviceType,
		}
	}

	if w.Capability != nil {
		response.Capability = CapabilityDTO{
			CapabilityID:   w.Capability.ID,
			CapabilityType: w.Capability.CapabilityType,
			ControlType:    w.Capability.ControlType,
		}
	}

	return response
}

type WidgetResponse struct {
	WidgetID     uint   `json:"widget_id"`
	WidgetOrder  uint   `json:"widget_order"`
	WidgetStatus string `json:"widget_status"`
	Value        uint   `json:"value"`

	Device     DeviceDTO     `json:"device"`
	Capability CapabilityDTO `json:"capability"`
}

type DeviceDTO struct {
	DeviceID   string   `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceLastHeartbeat time.Time `json:"device_last_heartbeat"`
	DeviceType string `json:"device_type"`
}

type CapabilityDTO struct {
	CapabilityID   uint   `json:"capability_id"`
	CapabilityType string `json:"capability_type"`
	ControlType    string `json:"control_type"`
}
