package main

import (
	"fmt"
	"project-home-iot/internal/auth"
	"project-home-iot/internal/database"

	"project-home-iot/internal/core/domain"
	"project-home-iot/internal/infrastructure/gorm"
	"project-home-iot/internal/infrastructure/mqtt"

	mqttlib "github.com/eclipse/paho.mqtt.golang"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	db := database.ConnectDB()
	db.AutoMigrate(
		&auth.UserAccount{},
		&domain.Device{},
		&domain.MonitorData{},
		&domain.Capability{},
		&domain.Widget{},
	)

	deviceRepo := gorm.NewDeviceRepository(db)
	mqttHandler := mqtt.NewMQTTHandler(deviceRepo)

	opts := mqttlib.NewClientOptions().AddBroker("tcp://localhost:1883")
	client := mqttlib.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 4. เริ่มต้น Subscribe ข้อมูล
	mqttHandler.SubscribeDeviceRegistration(client)
	mqttHandler.SubscribeSensorData(client)

	handler := InitializeApp(db)

	app := fiber.New()

	SetupRoutes(app, handler)

	app.Listen(":3000")
}
