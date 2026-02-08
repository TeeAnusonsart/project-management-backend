package domain

type Recorder interface {
	DataReceive(log *Log) error
	RecordLog(widgetId uint,eventType string,value uint) error
}