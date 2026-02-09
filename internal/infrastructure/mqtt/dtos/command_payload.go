package dtos

type DeviceCommand struct {
	CapabilityID       string `json:"capability_id"`
	Value         uint   `json:"value"`
	ReplyTopic    string `json:"reply_topic"`
}

type CommandResponse struct {
	Status        string `json:"status"`
}