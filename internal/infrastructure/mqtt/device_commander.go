// internal/infrastructure/mqtt/device_commander.go
package mqtt

import (
	"encoding/json"
	"fmt"
	"project-home-iot/internal/core/domain"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTDeviceCommander struct {
	client mqtt.Client
}

func NewMQTTDeviceCommander(client mqtt.Client) domain.DeviceCommander {
	return &MQTTDeviceCommander{client: client}
}

func (m *MQTTDeviceCommander) RequestCommand(deviceId string, cmd *domain.DeviceCommand, correlationID string) (*domain.CommandResponse, error) {
	replyTopic := fmt.Sprintf("devices/%s/command/response", correlationID)

	ch := make(chan domain.CommandResponse, 1)

	tokenSub := m.client.Subscribe(replyTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
		var respDTO CommandResponse
		if err := json.Unmarshal(msg.Payload(), &respDTO); err == nil {
			ch <- domain.CommandResponse{
				Status: respDTO.Status,
			}
		}
		c.Unsubscribe(replyTopic)
	})

	if tokenSub.Wait() && tokenSub.Error() != nil {
		return nil, tokenSub.Error()
	}

	payloadDTO := DeviceCommandDomainToPayload(cmd)
	payload, err := json.Marshal(payloadDTO)
	if err != nil {
		fmt.Printf("Error marshalling command payload: %v\n", err)
		return nil, err
	}
	topic := fmt.Sprintf("devices/%s/command/request", deviceId)
	tokenPub := m.client.Publish(topic, 1, false, payload)
	if tokenPub.Wait() && tokenPub.Error() != nil {
		m.client.Unsubscribe(replyTopic)
		return nil, tokenPub.Error()
	}

	if err := tokenPub.Error(); err != nil {
		return nil, err
	}


	select {
	case res := <-ch:
		return &res, nil
	case <-time.After(60 * time.Second):
		m.client.Unsubscribe(replyTopic)
		return nil, fmt.Errorf("device timeout on topic %s", replyTopic)
	}
}

func DeviceCommandDomainToPayload(d *domain.DeviceCommand) *DeviceCommand {
	return &DeviceCommand{
		CapabilityType: d.CapabilityType,
		ControlType: d.ControlType,
		Value:        d.Value,
		ReplyTopic:   d.ReplyTopic,
	}
}

type DeviceCommand struct {
	CapabilityType string `json:"capability_type"`
	ControlType    string `json:"control_type"`
	Value          string   `json:"value"`
	ReplyTopic     string `json:"reply_topic"`
}

type CommandResponse struct {
	Status string `json:"status"`
}

