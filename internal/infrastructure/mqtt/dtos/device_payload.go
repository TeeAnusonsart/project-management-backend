package dto

type DevicePayload struct {
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	DeviceKey  string `json:"device_key"`
	Topic      string `json:"topic"`
}
