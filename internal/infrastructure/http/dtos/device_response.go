package dtos

import "time"

type DeviceResponse struct {
	DeviceID         string   `json:"device_id"`
	DeviceLastHeartbeat time.Time   `json:"device_last_heartbeat"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
}
