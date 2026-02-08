package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

type Subscriber struct {
	client mqtt.Client
}

func NewSubscriber(client mqtt.Client) *Subscriber {
	return &Subscriber{client: client}
}

func (s *Subscriber) SubscribeSensor(handler mqtt.MessageHandler) error {
	token := s.client.Subscribe(
		"/device/+/sensor",
		1,
		handler,
	)
	token.Wait()
	return token.Error()
}
