package domain

type WidgetRepository interface {
	CreateWidget(widget *Widget) error
	UpdateValue(widgetID uint, value uint) error
	// FindByID(widgetID uint) (*Widget, error)
	// CreateMonitorData(data *MonitorData) error
}
