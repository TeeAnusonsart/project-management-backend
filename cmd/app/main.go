package main

import (
	"fmt"
	"project-home-iot/internal/database"

	// "project-home-iot/internal/core/domain"
	"project-home-iot/internal/core/usecase"
	"project-home-iot/internal/infrastructure/gorm"
	"project-home-iot/internal/infrastructure/gorm/models"
	"project-home-iot/internal/infrastructure/http"
	"project-home-iot/internal/infrastructure/mqtt"
	// "project-home-iot/internal/interfaces/http"

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
	db.AutoMigrate(&models.Log{})

	// 2. สร้างตารางลูกที่มี FK ทีหลัง
	err := db.AutoMigrate(&models.Widget{})

	if err != nil {
		fmt.Printf("Migration Error: %v\n", err) // ดูว่ามันด่าว่าอะไร
	}

	deviceRepo := gorm.NewDeviceRepository(db)
	capabilityRepo := gorm.NewCapabilityRepository(db)
	widgetRepo := gorm.NewWidgetRepository(db)
	recorderRepo := gorm.NewRecorderRepository(db)

	// capabilityUsecase := usecase.NewCapabilityUsecase(capabilityRepo)
	widgetUsecase := usecase.NewWidgetUsecase(widgetRepo)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo, capabilityRepo, widgetRepo)
	mqttHandler := mqtt.NewMQTTHandler(deviceUsecase,widgetUsecase)

	opts := mqttlib.NewClientOptions().AddBroker("tcp://localhost:1883")
	client := mqttlib.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 4. เริ่มต้น Subscribe ข้อมูล
	mqttHandler.SubscribeDeviceRegistration(client)
	// mqttHandler.SubscribeReply(client)
	// mqttHandler.SubscribeSensorData(client)

	handler := InitializeApp(db)

	app := fiber.New()

	deviceCommander :=mqtt.NewMQTTDeviceCommander(client)
	commandUsecase := usecase.NewCommandUsecase(deviceRepo,widgetRepo,deviceCommander,recorderRepo)

	commandHandler := http.NewCommandHandler(commandUsecase)
	app.Post("/api/widgets/:widgetId/command", commandHandler.SendCommand)


	SetupRoutes(app, handler)
	fmt.Println("Server is starting on :3000...")
	app.Listen(":3000")
	// if err := app.Listen(":3000"); err != nil {
	// 	fmt.Printf("Error starting server: %v\n", err)
	// }
}
