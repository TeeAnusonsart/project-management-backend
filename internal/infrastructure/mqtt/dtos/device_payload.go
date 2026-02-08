package dtos

type DevicePayload struct {
	DeviceID   uint   `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	Topic      string `json:"topic"`
}

