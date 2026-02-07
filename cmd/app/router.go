package main

import (
	"project-home-iot/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, handlers *AppHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.Post("/login", handlers.Auth.Login)
	authGroup.Post("/register", handlers.Auth.Register)

	secretGroup := api.Group("/secret", middleware.JWTMiddleware())
	secretGroup.Get("/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"secret_data": "นี่คือข้อมูลลับที่ได้รับการป้องกันด้วย JWT",
		})
	})

	deviceGroup := api.Group("/devices")
	deviceGroup.Get("/", handlers.Device.ListDevices)
	deviceGroup.Get("/:device_id", handlers.Device.GetDevice)
	deviceGroup.Put("/:device_id", handlers.Device.UpdateDevice)
	deviceGroup.Post("/:device_id/pair", handlers.Device.PairDevice)
	deviceGroup.Post("/:device_id/unpair", handlers.Device.UnpairDevice)

	api.Post("/widgets/:widgetId/command", handlers.Command.SendCommand)
}
