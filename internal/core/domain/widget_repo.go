package domain

type WidgetRepository interface {
	GetWidgetIdByDeviceAndCapability(deviceId uint, capabilityId uint) *Widget
	CreateWidget(widget *Widget) error
	Update(widget *Widget) error
	Delete(id uint) error

	FindAll() ([]*Widget, error)
	FindByID(id uint) (*Widget, error)

	UpdateValue(widgetID uint, value uint) error
	UpdateStatus(id uint, status string) error
	ChangeOrder(roomID uint, widgetOrders []uint) error

}
