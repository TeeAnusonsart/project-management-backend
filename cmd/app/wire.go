//go:build wireinject
// +build wireinject

package main

import (
	"project-home-iot/internal/auth"
	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	gormRepo "project-home-iot/internal/infrastructure/gorm"
	httpHandler "project-home-iot/internal/infrastructure/http"
	mqttInfra "project-home-iot/internal/infrastructure/mqtt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type AppHandlers struct {
	Auth    *auth.AuthHandler
	Device  *httpHandler.DeviceHandler
	Command *httpHandler.CommandHandler
}

func InitializeApp(db *gorm.DB, client mqtt.Client) *AppHandlers {
	wire.Build(
    auth.ProviderSet,

    gormRepo.NewDeviceRepository,
    gormRepo.NewCapabilityRepository,
    gormRepo.NewWidgetRepository,

    wire.Bind(new(domain.DeviceRepository), new(*gormRepo.DeviceRepository)),
    wire.Bind(new(domain.CapabilityRepository), new(*gormRepo.CapabilityRepository)),
    wire.Bind(new(domain.WidgetRepository), new(*gormRepo.WidgetRepository)),

    usecase.NewCapabilityUsecase,
    usecase.NewDeviceUsecase,
    usecase.NewWidgetUsecase,

    mqttInfra.NewMQTTDeviceCommander,

    usecase.NewCommandUsecase,

    httpHandler.NewDeviceHandler,
    httpHandler.NewCommandHandler,

    wire.Struct(new(AppHandlers), "*"),
)


	return &AppHandlers{}
}
