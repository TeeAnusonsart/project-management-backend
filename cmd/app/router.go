package main

import (
	"project-home-iot/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handlers *AppHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	authGroup := app.Group("/api/auth")
	authGroup.Post("/login", handlers.Auth.Login)
	authGroup.Post("/register", handlers.Auth.Register)

	secretGroup := app.Group("/api/secret", middleware.JWTMiddleware())
	secretGroup.Get("/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"secret_data": "นี่คือข้อมูลลับที่ได้รับการป้องกันด้วย JWT",
		})
	})

}
