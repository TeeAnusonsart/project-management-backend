package usecase

// ในไฟล์ usecase/device_map.go (หรือที่เก็บ Map เดิม)

type CapabilityRef struct {
    Type    string
    Control string
}

var DeviceCapabilityMap = map[string][]CapabilityRef{
    "light": {
        {Type: "toggle", Control: "light-switch"},
		{Type: "mode", Control: "light-system"},
		{Type: "sensor", Control: "light-intensity"},
    },
    "air": {
        {Type: "toggle", Control: "air-power"},
        {Type: "adjust", Control: "atemperature"},
    },
    "fan": {
        {Type: "toggle", Control: "fan-power"},
        {Type: "speed", Control: "dropdown"},
    },
	"air-samsung": {
        {Type: "toggle", Control: "air-power"},
		{Type: "toggle", Control: "air-bomb"},
        {Type: "adjust", Control: "temperature"},
    },
}
