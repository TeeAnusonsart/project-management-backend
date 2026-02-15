package mqtt

import (
	"encoding/json"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"project-home-iot/internal/core/usecase"
	dtos "project-home-iot/internal/infrastructure/mqtt/dtos"
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

	// deviceID, err := strconv.ParseUint(deviceIDStr, 10, 64)
	// if err != nil {
	// 	return
	// }

	var payload dtos.SensorPayload
	if err := json.Unmarshal(m.Payload(), &payload); err != nil {
		return
	}

	if payload.CapabilityType == "heartbeat" { 
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
	)
}
