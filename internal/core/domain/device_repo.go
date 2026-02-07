package domain

type DeviceRepository interface {
	CreateDevice(device *Device) error
	FindByWidgetID(widgetID uint) (*Device, error)
	// CreateMonitorData(data *MonitorData) error
}
