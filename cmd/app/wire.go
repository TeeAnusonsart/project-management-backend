//go:build wireinject
// +build wireinject

package main

import (
	"project-home-iot/internal/core/domain"
	"firebase.google.com/go/v4/auth"
	"project-home-iot/internal/core/usecase"
	gormRepo "project-home-iot/internal/infrastructure/gorm"
	httpHandler "project-home-iot/internal/infrastructure/http"
	mqttInfra "project-home-iot/internal/infrastructure/mqtt"
	"project-home-iot/internal/firebase"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/wire"
	"gorm.io/gorm"
)

type AppHandlers struct {
	// Auth    *auth.AuthHandler
	AuthClient *auth.Client
	Device  *httpHandler.DeviceHandler
	Command *httpHandler.CommandHandler
    Room    *httpHandler.RoomHandler
    Widget  *httpHandler.WidgetHandler
    User *httpHandler.UserHandler
    MQTT *mqttInfra.MQTTHandler
	SensorSubscriber *mqttInfra.SensorSubscriber
	SensorHandler *mqttInfra.SensorHandler
	UserRepo domain.UserRepository
}

func InitializeApp(db *gorm.DB, client mqtt.Client) (*AppHandlers, error) {
	wire.Build(
		// auth
		// auth.ProviderSet,

		// repo
		gormRepo.NewDeviceRepository,
		gormRepo.NewCapabilityRepository,
		gormRepo.NewWidgetRepository,
		gormRepo.NewRoomRepository,
        gormRepo.NewUserRepository,
		gormRepo.NewRecorderRepository,
	

		// bind
		wire.Bind(new(domain.DeviceRepository), new(*gormRepo.DeviceRepository)),
		wire.Bind(new(domain.CapabilityRepository), new(*gormRepo.CapabilityRepository)),
		wire.Bind(new(domain.WidgetRepository), new(*gormRepo.WidgetRepository)),
		wire.Bind(new(domain.RoomRepository), new(*gormRepo.RoomRepository)),
        wire.Bind(new(domain.UserRepository), new(*gormRepo.UserRepository)),
		wire.Bind(new(domain.Recorder), new(*gormRepo.RecorderRepository)),

		firebase.NewFirebaseAuthClient,

		// uc
		usecase.NewDeviceUsecase,
		usecase.NewWidgetUsecase,
		usecase.NewRoomUsecase,
        usecase.NewUserUsecase,
        usecase.NewCommandUsecase,
		usecase.NewRecordLogUsecase,
        
		// mqtt
		mqttInfra.NewMQTTDeviceCommander,
        mqttInfra.NewMQTTPairCommander,
		mqttInfra.NewSensorSubscriber,
		mqttInfra.NewSensorHandler,



		// handler
		httpHandler.NewDeviceHandler,
		httpHandler.NewCommandHandler,
		httpHandler.NewRoomHandler,
        httpHandler.NewWidgetHandler,
        httpHandler.NewUserHandler,
        mqttInfra.NewMQTTHandler,

		wire.Struct(new(AppHandlers), "*"),
	)

	return &AppHandlers{},nil
}
