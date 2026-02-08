package domain

type Widget struct {
	ID uint

	WidgetOrder  uint
	WidgetStatus string
	Value        uint

	DeviceID     uint
	CapabilityID uint

	Device     *Device
	Capability *Capability
}
