package usecase

import "project-home-iot/internal/core/domain"

type WidgetUsecase interface {
	CreateWidget(widget *domain.Widget) error
	UpdateWidget(widget *domain.Widget) error
	DeleteWidget(id uint) error

	ListWidgets() ([]*domain.Widget, error)
	GetWidget(id uint) (*domain.Widget, error)
	GetWidgetByStatus(status string) ([]*domain.Widget, error)
	GetLogs(id uint)([]*domain.Log, error)
	ListWidgetsByRoom(roomID uint) ([]*domain.Widget, error)

	UpdateValue(widgetID uint, value string) error
	UpdateStatus(id uint, status string) error
	ChangeOrder(roomID uint, widgetOrders []uint) error
}

type widgetUsecase struct {
	widgetRepo domain.WidgetRepository
	recorderRepo domain.Recorder
}

func NewWidgetUsecase(wr domain.WidgetRepository,rr domain.Recorder) WidgetUsecase {
	return &widgetUsecase{widgetRepo: wr,recorderRepo: rr}
}

func (u *widgetUsecase) CreateWidget(widget *domain.Widget) error {
	return u.widgetRepo.CreateWidget(widget)
}

func (u *widgetUsecase) UpdateWidget(widget *domain.Widget) error {
	return u.widgetRepo.Update(widget)
}

func (u *widgetUsecase) DeleteWidget(id uint) error {
	return u.widgetRepo.Delete(id)
}

func (u *widgetUsecase) GetWidgetByStatus(status string) ([]*domain.Widget, error) {
	return u.widgetRepo.GetWidgetByStatus(status)
}

func (u *widgetUsecase) ListWidgets() ([]*domain.Widget, error) {
	return u.widgetRepo.FindAll()
}

func (u *widgetUsecase) GetWidget(id uint) (*domain.Widget, error) {
	return u.widgetRepo.FindByID(id)
}

func (u *widgetUsecase) GetLogs(id uint) ([]*domain.Log, error){
	return u.recorderRepo.GetLogByWidgetID(id)
}

func (u *widgetUsecase) ListWidgetsByRoom(roomID uint) ([]*domain.Widget, error) {
	return u.widgetRepo.FindByRoomID(roomID)
}

func (u *widgetUsecase) UpdateValue(widgetID uint, value string) error {
	return u.widgetRepo.UpdateValue(widgetID, value)
}

func (u *widgetUsecase) UpdateStatus(id uint, status string) error {
	return u.widgetRepo.UpdateStatus(id, status)
}

func (u *widgetUsecase) ChangeOrder(roomID uint, widgetOrders []uint) error {
    return u.widgetRepo.ChangeOrder(roomID, widgetOrders)
}

