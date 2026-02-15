package dtos

import "time"


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
}
