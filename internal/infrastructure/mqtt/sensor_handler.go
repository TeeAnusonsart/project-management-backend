package mqtt

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"project-home-iot/internal/core/usecase"
)

type SensorHandler struct {
	recordLogUC usecase.RecordLogUsecase
	deviceUC usecase.DeviceUsecase
}

func NewSensorHandler(uc usecase.RecordLogUsecase,duc usecase.DeviceUsecase) *SensorHandler {
	return &SensorHandler{
		recordLogUC: uc,
		deviceUC: duc,
	}
}

func (h *SensorHandler) HandleSensorMessage(
	c mqtt.Client,
	m mqtt.Message,
) {
	topic := m.Topic()
	parts := strings.Split(topic, "/")

	deviceID := parts[1]

	var payload SensorPayload
	if err := json.Unmarshal(m.Payload(), &payload); err != nil {
		return
	}

	if payload.CapabilityType == "status" && payload.ControlType == "heartbeat" { 
		_ = h.deviceUC.UpdateHeartbeat(deviceID)
		return
	}

	_ = h.recordLogUC.Execute(
		deviceID,
		payload.CapabilityType,
		payload.ControlType,
		// payload.CapabilityID,
		"sensor",
		payload.Value,
		payload.Actor,

	)
}

type SensorPayload struct {
	// CapabilityID uint `json:"capability_id"`
	CapabilityType string `json:"capability_type"`
	ControlType string `json:"control_type"`
	Value string `json:"value"`
	Actor string `json:"actor"`
}
