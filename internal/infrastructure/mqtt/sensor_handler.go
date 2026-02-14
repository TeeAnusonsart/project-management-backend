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
	topic := m.Topic()
	parts := strings.Split(topic, "/")

	deviceIDStr := parts[2]

	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 64)
	if err != nil {
		return
	}

	var payload dtos.SensorPayload
	if err := json.Unmarshal(m.Payload(), &payload); err != nil {
		return
	}

	_ = h.recordLogUC.Execute(
		uint(deviceID),
		payload.CapabilityID,
		"sensor",
		payload.Value,
	)
}
