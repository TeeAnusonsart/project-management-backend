package main

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handlers *AppHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	authGroup := app.Group("/api/auth")
	authGroup.Post("/login", handlers.Auth.Login)
	authGroup.Post("/register", handlers.Auth.Register)

}
