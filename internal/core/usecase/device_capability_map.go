package usecase

var DeviceCapabilityMap = map[string][]string{
	"light": {"switch"},
	"air":   {"power", "temperature"},
	"fan":   {"power", "speed"},
}
