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
    Room    *httpHandler.RoomHandler
    Widget  *httpHandler.WidgetHandler
    User *httpHandler.UserHandler
}

func InitializeApp(db *gorm.DB, client mqtt.Client) *AppHandlers {
	wire.Build(
		// auth
		auth.ProviderSet,

		// repo
		gormRepo.NewDeviceRepository,
		gormRepo.NewCapabilityRepository,
		gormRepo.NewWidgetRepository,
		gormRepo.NewRoomRepository,
        gormRepo.NewUserRepository,

		// bind
		wire.Bind(new(domain.DeviceRepository), new(*gormRepo.DeviceRepository)),
		wire.Bind(new(domain.CapabilityRepository), new(*gormRepo.CapabilityRepository)),
		wire.Bind(new(domain.WidgetRepository), new(*gormRepo.WidgetRepository)),
		wire.Bind(new(domain.RoomRepository), new(*gormRepo.RoomRepository)),
        wire.Bind(new(domain.UserRepository), new(*gormRepo.UserRepository)),

		// uc
		usecase.NewDeviceUsecase,
		usecase.NewWidgetUsecase,
		usecase.NewRoomUsecase,
        usecase.NewUserUsecase,

		// mqtt
		mqttInfra.NewMQTTDeviceCommander,
		usecase.NewCommandUsecase,

		// handler
		httpHandler.NewDeviceHandler,
		httpHandler.NewCommandHandler,
		httpHandler.NewRoomHandler,
        httpHandler.NewWidgetHandler,
        httpHandler.NewUserHandler,

		wire.Struct(new(AppHandlers), "*"),
	)

	return &AppHandlers{}
}
