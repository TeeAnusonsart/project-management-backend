package dtos

type UpdateDeviceRequest struct {
	DeviceName string `json:"device_name"`
}

type PairDeviceRequest struct {
	DeviceKey string `json:"device_key"`
}
