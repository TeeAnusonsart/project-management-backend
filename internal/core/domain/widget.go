package domain

import "time"

type Widget struct {
	ID uint

	WidgetOrder  uint
	WidgetStatus string
	Value        string

	DeviceID     string
	CapabilityID uint

	Device     *Device
	Capability *Capability
}

type Log struct {
	ID        uint
	WidgetID  uint
	Actor     string
	Value     string
	EventType string
	CreatedAt  time.Time
}

type WidgetRepository interface {
	GetWidgetIdByDeviceAndCapability(deviceId string, capabilityId uint) *Widget
	CreateWidget(widget *Widget) error
	Update(widget *Widget) error
	Delete(id uint) error

	DeleteByDeviceId(deviceID string) error

	FindAll() ([]*Widget, error)
	FindByID(id uint) (*Widget, error)
	FindByRoomID(roomID uint) ([]*Widget, error)
	GetWidgetByStatus(status string) ([]*Widget, error)
    FindByRoomWithStatus(roomID uint, status string) ([]*Widget, error)

	UpdateValue(widgetID uint, value string) error
	UpdateStatus(id uint, status string) error
	ChangeOrder(roomID uint, widgetOrders []uint) error
}

type Recorder interface {
	// DataReceive(log *Log) error
	// RecordLog(widgetId uint,eventType string,value uint) error
	RecordLog(log *Log) error
	GetLogByWidgetID(widgetID uint) ([]*Log, error)
}
