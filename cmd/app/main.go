package main

import (
	"errors"
	"fmt"
	// "project-home-iot/internal/auth"
	"log"
	"project-home-iot/internal/database"
	"project-home-iot/internal/infrastructure/gorm"
	"project-home-iot/internal/core/domain"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "project-home-iot/cmd/app/docs"
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
		// &auth.UserAccount{},
	)

	gorm.SeedAll(db)

	// db.AutoMigrate(&models.Widget{})

	// opts := mqttlib.NewClientOptions().AddBroker("tcp://mqtt-broker:1883").SetClientID("go-backend-server")

	opts := mqttlib.NewClientOptions().AddBroker("tcp://192.168.137.241:1883").SetClientID("go-backend-server")
	// opts := mqttlib.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-backend-server")

	// opts := mqttlib.NewClientOptions().
	// AddBroker("tcp://broker.hivemq.com:1883").
	// SetClientID("go-backend-server")
	client := mqttlib.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// handlers := InitializeApp(db, client)

	handlers, err := InitializeApp(db, client)
	if err != nil {
		log.Fatalf("ไม่สามารถเริ่มระบบได้: %v", err)
	}

	handlers.MQTT.SubscribeDeviceRegistration(client)
	client.Publish("xxx/y", 0, false, "Hello MQTT")
	_ = handlers.SensorSubscriber.SubscribeSensor(handlers.SensorHandler.HandleSensorMessage)
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				message = fiberErr.Message
			}

			return c.Status(code).JSON(domain.ErrorResponse{
				Status:  "error",
				Code:    code,
				Message: message,
				Errors:  nil,
			})
		},
	})

	app.Static("/uploads", "./uploads")
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	SetupRoutes(app, handlers)

	fmt.Println("Server is starting on :3000...")
	app.Listen(":3000")
}
