package domain

type WidgetRepository interface {
	CreateWidget(widget *Widget) error
	// CreateMonitorData(data *MonitorData) error
}
