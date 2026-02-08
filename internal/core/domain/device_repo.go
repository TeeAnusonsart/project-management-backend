package domain

type DeviceRepository interface {
	CreateDevice(device *Device) error
	FindByWidgetID(widgetID uint) (*Device, error)

	GetAllSummaries() ([]*DeviceSummary, error)
	GetSummaryByID(id uint) (*DeviceSummary, error)

	GetByID(id uint) (*Device, error)

	UpdateName(id uint, name string) error
	Pair(id uint, deviceKey string) error
	Unpair(id uint) error
}

