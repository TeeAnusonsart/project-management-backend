package usecase

var DeviceCapabilityMap = map[string][]string{
	"light": {"toggle"},
	"air":   {"toggle", "adjust"},
	"fan":   {"toggle", "adjust"},
}
