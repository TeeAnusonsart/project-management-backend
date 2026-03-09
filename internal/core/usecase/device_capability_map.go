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
        {Type: "adjust", Control: "temperature"},
    },
    "fan": {
        {Type: "toggle", Control: "fan-power"},
        {Type: "level", Control: "fan-speed"},
    },
	"temperature-sensor": {
        {Type: "sensor", Control: "temperature"},
    },
    "temperature_humidity": {
        {Type: "sensor", Control: "humidity"},
        {Type: "adjust", Control: "content"},
        {Type: "mode", Control: "display"},
    },
    "environment": {
        {Type: "sensor", Control: "humidity"},
        {Type: "sensor", Control: "gas"},
        {Type: "text", Control: "content"},

    },
    "ph": {
        {Type: "sensor", Control: "ph"},
    },
}
