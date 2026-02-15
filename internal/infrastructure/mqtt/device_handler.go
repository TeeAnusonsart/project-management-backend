// internal/infrastructure/mqtt/device_handler.go
package mqtt

import (
	"encoding/json"

	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
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

	client.Subscribe(topic, 1, func(c mqtt.Client, m mqtt.Message) {

		var payload DevicePayload

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			return
		}

		device := PayloadToDomain(&payload)

		
		err = h.deviceUsecase.RegisterDevice(device)
		if err != nil {
			return
		}

	})
}


func (h *MQTTHandler) DeviceCommand(client mqtt.Client) {

}

type DevicePayload struct {
	DeviceID   string   `json:"device_id"`
	DeviceName string `json:"device_name"`
	DeviceType string `json:"device_type"`
	Topic      string `json:"topic"`
}

func PayloadToDomain(p *DevicePayload) *domain.Device {
	return &domain.Device{
		DeviceID:         p.DeviceID,
		DeviceName: p.DeviceName,
		DeviceType: p.DeviceType,
	}
}