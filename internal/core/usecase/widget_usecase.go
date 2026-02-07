package usecase

import (
	"project-home-iot/internal/core/domain"
)

type WidgetUsecase interface {
	CreateWidget(widget *domain.Widget) error
	UpdateValue(widgetID uint, value uint) error
}

type widgetUsecase struct {
	widgetRepo domain.WidgetRepository
}

func NewWidgetUsecase(wr domain.WidgetRepository) WidgetUsecase {
	return &widgetUsecase{widgetRepo: wr}
}

func (u *widgetUsecase) CreateWidget(widget *domain.Widget) error {
	return u.widgetRepo.CreateWidget(widget)
}

func (u *widgetUsecase) UpdateValue(widgetID uint, value uint) error {
	return u.widgetRepo.UpdateValue(widgetID, value)
}
