package domain

type DeviceRepository interface {
	CreateDevice(device *Device) error
	// CreateMonitorData(data *MonitorData) error
}
