package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"project-home-iot/internal/core/domain"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTPairCommander struct {
	client mqtt.Client
}

func NewMQTTPairCommander(client mqtt.Client) domain.PairCommander {
	return &MQTTPairCommander{client: client}
}

func (m *MQTTPairCommander) RequestPair(deviceID string, deviceKey string) error {
	requestTopic := fmt.Sprintf("devices/%s/pair/request", deviceID)
	replyTopic := fmt.Sprintf("devices/%s/pair/response", deviceID)
	
	ch := make(chan string, 1)

	tokenSub := m.client.Subscribe(replyTopic, 1, func(c mqtt.Client, msg mqtt.Message) {
		var resp struct {
			Status string `json:"status"`
		}

		if err := json.Unmarshal(msg.Payload(), &resp); err == nil {
			ch <- resp.Status
		}

		c.Unsubscribe(replyTopic)
	})

	if tokenSub.Wait() && tokenSub.Error() != nil {
		return tokenSub.Error()
	}

	payload := map[string]interface{}{
		"device_key": deviceKey,
	}

	data, _ := json.Marshal(payload)

	tokenPub := m.client.Publish(requestTopic, 1, false, data)
	if tokenPub.Wait() && tokenPub.Error() != nil {
		m.client.Unsubscribe(replyTopic)
		return tokenPub.Error()
	}

	select {
	case status := <-ch:
		if status != "success" {
			return fmt.Errorf("pair failed")
		}
		return nil

	case <-time.After(30 * time.Second):
		m.client.Unsubscribe(replyTopic)
		return fmt.Errorf("pair timeout")
	}
}

func (m *MQTTPairCommander) Subscribe(deviceID string) error {
	topic := fmt.Sprintf("devices/%s/pair/response", deviceID)
	token := m.client.Subscribe(topic, 1, nil)
	token.Wait()
	return token.Error()
}
