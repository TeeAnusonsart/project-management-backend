package dtos

type DevicePayload struct {
	DeviceID   uint   `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	Topic      string `json:"topic"`
}

type DeviceCommand struct {
	CapabilityType string `json:"capability_type"`
	ControlType    string `json:"control_type"`
	Value          uint   `json:"value"`
	ReplyTopic     string `json:"reply_topic"`
}

type CommandResponse struct {
	Status string `json:"status"`
}
