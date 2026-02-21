package main

import (
	"fmt"
	"project-home-iot/internal/auth"
	"project-home-iot/internal/database"
	"project-home-iot/internal/infrastructure/gorm"
	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	// "project-home-iot/internal/core/usecase"
	// "project-home-iot/internal/infrastructure/gorm"
	// "project-home-iot/internal/infrastructure/mqtt"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	db := database.ConnectDB()

	fmt.Println("Starting fresh migration...")
	db.AutoMigrate(
		&gorm.User{},
		&gorm.Room{},
		&gorm.Capability{},
		&gorm.Device{}, 
		&gorm.Widget{}, 
		&gorm.Log{},
		&auth.UserAccount{},
	)

	gorm.SeedAll(db)

	// db.AutoMigrate(&models.Widget{})

	// opts := mqttlib.NewClientOptions().AddBroker("tcp://mqtt-broker:1883").SetClientID("go-backend-server")
	
	opts := mqttlib.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-backend-server")

	client := mqttlib.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	handlers := InitializeApp(db, client)

	// recorderRepo := gorm.NewRecorderRepository(db)
	// widgetRepo := gorm.NewWidgetRepository(db)
	// recordLogUC := usecase.NewRecordLogUsecase(recorderRepo,widgetRepo)

	// handler := mqtt.NewSensorHandler(recordLogUC)
	// sub := mqtt.NewSubscriber(client)

	// _ = sub.SubscribeSensor(handler.HandleSensorMessage)

	handlers.MQTT.SubscribeDeviceRegistration(client)
	client.Publish("xxx/y", 0, false, "Hello MQTT")
	_ = handlers.SensorSubscriber.SubscribeSensor(handlers.SensorHandler.HandleSensorMessage)
	app := fiber.New()

	app.Static("/uploads", "./uploads")

	SetupRoutes(app, handlers)

	fmt.Println("Server is starting on :3000...")
	app.Listen(":3000")
}
