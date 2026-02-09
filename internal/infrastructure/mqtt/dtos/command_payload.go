package dtos

type DeviceCommand struct {
	WidgetID      uint   `json:"widget_id"`
	CapabilityID       string `json:"capability_id"`
	Value         uint   `json:"value"`
	ReplyTopic    string `json:"reply_topic"`
	CorrelationID string `json:"correlation_id"`
}

type CommandResponse struct {
	Status        string `json:"status"`
}