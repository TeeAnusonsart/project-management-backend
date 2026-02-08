package dtos

type DeviceResponse struct {
	ID         uint   `json:"id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
}
