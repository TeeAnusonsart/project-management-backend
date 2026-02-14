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

type Log struct {
	ID        uint
	WidgetID  uint
	ActorType string
	Value     uint
	EventType string
}

type WidgetRepository interface {
	GetWidgetIdByDeviceAndCapability(deviceId uint, capabilityId uint) *Widget
	CreateWidget(widget *Widget) error
	Update(widget *Widget) error
	Delete(id uint) error

	FindAll() ([]*Widget, error)
	FindByID(id uint) (*Widget, error)
	FindByRoomID(roomID uint) ([]*Widget, error)

	UpdateValue(widgetID uint, value uint) error
	UpdateStatus(id uint, status string) error
	ChangeOrder(roomID uint, widgetOrders []uint) error

}

type Recorder interface {
	// DataReceive(log *Log) error
	RecordLog(widgetId uint,eventType string,value uint) error
}
