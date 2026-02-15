// internal/infrastructure/mqtt/device_handler.go
package mqtt

import (
	"encoding/json"
	"fmt"

	// "project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"project-home-iot/internal/infrastructure/mappers"
	dto "project-home-iot/internal/infrastructure/mqtt/dtos"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	deviceUsecase usecase.DeviceUsecase
	widgetUsecase usecase.WidgetUsecase
}

func NewMQTTHandler(uc usecase.DeviceUsecase,wuc usecase.WidgetUsecase) *MQTTHandler {
	return &MQTTHandler{deviceUsecase: uc,widgetUsecase: wuc}
}

func (h *MQTTHandler) SubscribeDeviceRegistration(client mqtt.Client) {
	topic := "devices/register"
	// fmt.Println("Subscribing to topic:", topic)

	client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
		// fmt.Println("Message received on topic:", m.Topic())
		// fmt.Println("Payload:", string(m.Payload()))
		var payload dto.DevicePayload

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			fmt.Printf("Error unmarshalling payload: %v\n", err)
			return
		}

		device := mappers.PayloadToDomain(&payload)

		
		err = h.deviceUsecase.RegisterDevice(device)
		if err != nil {
			// fmt.Printf("Failed to save device %s: %v\n", payload.DeviceID, err)
			return
		}

		// fmt.Printf("Successfully registered device: %s\n", payload.DeviceID)
	})
}

// func (h *MQTTHandler) SubscribeReply(client mqtt.Client) {

// 	client.Subscribe("devices/reply/+", 1, func(c mqtt.Client, m mqtt.Message) {

// 		var resp dto.CommandResponse
// 		json.Unmarshal(m.Payload(), &resp)

// 		fmt.Printf("Received reply for correlation ID %s: Widget %d new value %d\n", resp.CorrelationID, resp.WidgetID, resp.Value)
// 		h.widgetUsecase.UpdateValue(resp.WidgetID, resp.Value)
// 	})
// }

// type MonitorPayload struct {
// 	Value    uint `json:"value"`
// 	WidgetID uint `json:"widget_id"`
// }

// func (h *MQTTHandler) SubscribeSensorData(client mqtt.Client) {
// 	topic := "devices/data"
// 	client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {
// 		var payload MonitorPayload

// 		if err := json.Unmarshal(m.Payload(), &payload); err != nil {
// 			fmt.Printf("Error sensor payload: %v\n", err)
// 			return
// 		}

// 		sensorData := &domain.MonitorData{
// 			Value:    payload.Value,
// 			WidgetID: payload.WidgetID,
// 		}

// 		if err := h.usecase.RecordSensorData(sensorData); err != nil {
// 			fmt.Printf("Failed to save sensor data: %v\n", err)
// 			return
// 		}

// 		fmt.Printf("Sensor Recorded: Widget %d = %d\n", payload.WidgetID, payload.Value)
// 	})
// }

func (h *MQTTHandler) DeviceCommand(client mqtt.Client) {

}