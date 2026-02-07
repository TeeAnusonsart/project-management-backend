package usecase

import (
	"project-home-iot/internal/core/domain"
)

type WidgetUsecase interface {
	CreateWidget(widget *domain.Widget) error
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
