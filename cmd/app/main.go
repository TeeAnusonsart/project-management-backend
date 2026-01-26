package main

import (
	"fmt"
	"project-home-iot/internal/auth"

	"project-home-iot/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	db := database.ConnectDB()
	db.AutoMigrate(&auth.UserAccount{})

	handler := InitializeApp(db)

	app := fiber.New()

	SetupRoutes(app, handler)

	app.Listen(":3000")
}
