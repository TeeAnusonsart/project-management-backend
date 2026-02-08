package domain

type PairCommander interface {
	RequestPair(deviceID uint, deviceKey string) error
	Subscribe(topic string) error
}

