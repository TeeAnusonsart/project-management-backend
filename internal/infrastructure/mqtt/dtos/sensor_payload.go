package dtos

type SensorPayload struct {
	// CapabilityID uint `json:"capability_id"`
	CapabilityType string `json:"capability_type"`
	ControlType string `json:"control_type"`
	Value uint `json:"value"`
}