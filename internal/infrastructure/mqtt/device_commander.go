// internal/infrastructure/mqtt/device_commander.go
package mqtt

import (
	"encoding/json"
	"fmt"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/mappers"
	dto "project-home-iot/internal/infrastructure/mqtt/dtos"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTDeviceCommander struct {
	client mqtt.Client
}

func NewMQTTDeviceCommander(client mqtt.Client) domain.DeviceCommander {
	return &MQTTDeviceCommander{client: client}
}

func (m *MQTTDeviceCommander) RequestCommand(deviceId uint, cmd *domain.DeviceCommand, correlationID string) (*domain.CommandResponse, error) {
	replyTopic := fmt.Sprintf("devices/%s/response", correlationID)

	ch := make(chan domain.CommandResponse, 1)

	// 1. Subscribe รอคำตอบ
	tokenSub := m.client.Subscribe(replyTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
		var respDTO dto.CommandResponse
		if err := json.Unmarshal(msg.Payload(), &respDTO); err == nil {
			// Map DTO กลับเป็น Domain Model
			ch <- domain.CommandResponse{
				Status: respDTO.Status,
				// Status, Message เพิ่มเติมตามต้องการ
			}
		}
		c.Unsubscribe(replyTopic)
	})

	if tokenSub.Wait() && tokenSub.Error() != nil {
		return nil, tokenSub.Error()
	}

	// 2. Publish คำสั่งออกไป

	payloadDTO := mappers.DeviceCommandDomainToPayload(cmd)

	fmt.Printf("📤 XXXPayload: %s\n", payloadDTO.ReplyTopic)
	fmt.Println("Is connected:", m.client.IsConnected())
	// fmt.Printf("📡 Sending cXXXXommand to device ID %d with correlation ID %s\n", deviceId, correlationID)

	payload, err := json.Marshal(payloadDTO)
	if err != nil {
		fmt.Printf("Error marshalling command payload: %v\n", err)
		return nil, err
	}
	topic := fmt.Sprintf("devices/%d/commands", deviceId)
	fmt.Printf("📡 Sending command to topic: %s\n", topic)
	tokenPub := m.client.Publish(topic, 1, false, payload)
	if tokenPub.Wait() && tokenPub.Error() != nil {
		m.client.Unsubscribe(replyTopic)
		return nil, tokenPub.Error()
	}

	if err := tokenPub.Error(); err != nil {
		fmt.Println("❌ Publish error:", err)
		return nil, err
	}

    fmt.Printf("📤 Command published to topic: %s\n", topic)

	// 3. รอผลลัพธ์ด้วย Timeout
	select {
	case res := <-ch:
		return &res, nil
	case <-time.After(60 * time.Second):
		m.client.Unsubscribe(replyTopic)
		return nil, fmt.Errorf("device timeout on topic %s", replyTopic)
	}
}
