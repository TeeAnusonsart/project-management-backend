package dtos

type DeviceCommand struct {
	WidgetID      uint   `json:"widget_id"`
	Command       string `json:"command"`
	Value         uint   `json:"value"`
	ReplyTopic    string `json:"reply_topic"`
	CorrelationID string `json:"correlation_id"`
}

type CommandResponse struct {
	CorrelationID string `json:"correlation_id"`
	WidgetID      uint   `json:"widget_id"`
	Value		 uint   `json:"value"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}