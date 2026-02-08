package domain

type WidgetRepository interface {
	CreateWidget(widget *Widget) error
	Update(widget *Widget) error
	Delete(id uint) error

	FindAll() ([]*Widget, error)
	FindByID(id uint) (*Widget, error)

	UpdateValue(widgetID uint, value uint) error
	UpdateStatus(id uint, status string) error
	ChangeOrder(widgetOrders []uint) error
}
