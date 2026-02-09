package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

type SensorSubscriber struct {
	client mqtt.Client
}

func NewSensorSubscriber(client mqtt.Client) *SensorSubscriber {
	return &SensorSubscriber{client: client}
}

func (s *SensorSubscriber) SubscribeSensor(handler mqtt.MessageHandler) error {
	token := s.client.Subscribe(
		"/device/+/sensor",
		1,
		handler,
	)
	token.Wait()
	return token.Error()
}
