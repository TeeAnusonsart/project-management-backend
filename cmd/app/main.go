package main

import (
	"fmt"
	"project-home-iot/internal/database"
	"project-home-iot/internal/infrastructure/gorm/models"
	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	db := database.ConnectDB()
	db.AutoMigrate(&models.Device{})
	db.AutoMigrate(&models.Capability{})
	db.AutoMigrate(&models.Widget{})

	opts := mqttlib.NewClientOptions().AddBroker("tcp://mqtt-broker:1883")

	client := mqttlib.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	handlers := InitializeApp(db, client)

	app := fiber.New()

	SetupRoutes(app, handlers)

	fmt.Println("Server is starting on :3000...")
	app.Listen(":3000")
}
