package mqtt

import (
	"encoding/json"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"project-home-iot/internal/core/usecase"
	dtos "project-home-iot/internal/infrastructure/mqtt/dtos"
)

type SensorHandler struct {
	recordLogUC *usecase.RecordLogUsecase
}

func NewSensorHandler(uc *usecase.RecordLogUsecase) *SensorHandler {
	return &SensorHandler{
		recordLogUC: uc,
	}
}

func (h *SensorHandler) HandleSensorMessage(
	c mqtt.Client,
	m mqtt.Message,
) {
	// topic: /device/{deviceId}/sensor/{eventType}
	topic := m.Topic()
	parts := strings.Split(topic, "/")

	// extract routing info
	deviceIDStr := parts[2]

	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 64)
	if err != nil {
		return
	}

	// parse payload
	var payload dtos.SensorPayload
	if err := json.Unmarshal(m.Payload(), &payload); err != nil {
		return
	}

	// 🔥 call usecase (domain safe)
	_ = h.recordLogUC.Execute(
		uint(deviceID),
		payload.CapabilityID,
		"sensor",
		payload.Value,
	)
}
