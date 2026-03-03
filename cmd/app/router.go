package main

import (
	"project-home-iot/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// func SetupRoutes(app *fiber.App, handlers *AppHandlers) {
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendString("Hello, World!")
// 	})

// 	// api := app.Group("/api")

// 	// authen
// 	authGroup := app.Group("/auth")
// 	authGroup.Post("/login", handlers.User.Login)
// 	// authGroup.Post("/register", handlers.Auth.Register)

// 	// adminOnly := middleware.RoleGuard(handlers.UserRepo, "ADMIN")

	
// 	// authMid := middleware.FirebaseAuth(handlers.AuthClient)
// 	// api := app.Group("/api",authMid)
// 	api := app.Group("/api")

// 	secretGroup := api.Group("/secret", middleware.JWTMiddleware())
// 	secretGroup.Get("/data", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"secret_data": "นี่คือข้อมูลลับที่ได้รับการป้องกันด้วย JWT",
// 		})
// 	})

// 	// device
// 	deviceGroup := api.Group("/devices")
// 	deviceGroup.Get("/", handlers.Device.ListDevices)
// 	deviceGroup.Get("/:device_id", handlers.Device.GetDevice)
// 	deviceGroup.Put("/:device_id", handlers.Device.UpdateDevice)
// 	deviceGroup.Post("/:device_id/pair", handlers.Device.PairDevice)
// 	deviceGroup.Post("/:device_id/unpair", handlers.Device.UnpairDevice)

// 	// room
// 	roomGroup := api.Group("/rooms")
// 	roomGroup.Get("/", handlers.Room.ListRooms)
// 	roomGroup.Post("/", handlers.Room.CreateRoom)
// 	roomGroup.Get("/:room_id", handlers.Room.GetRoom)
// 	roomGroup.Put("/:room_id", handlers.Room.UpdateRoom)
// 	roomGroup.Delete("/:room_id", handlers.Room.DeleteRoom)
// 	roomGroup.Delete("/:room_id", handlers.Room.DeleteRoom)

// 	roomGroup.Post("/:room_id/devices", handlers.Room.AddDevice)
// 	roomGroup.Get("/:room_id/devices", handlers.Room.ListDevices)
// 	roomGroup.Post("/:room_id/widgets/order", handlers.Widget.ChangeOrder)

// 	roomGroup.Get("/:room_id/widgets", handlers.Widget.ListByRoom)

// 	// widget
// 	widgetGroup := api.Group("/widgets")

// 	widgetGroup.Get("/", handlers.Widget.ListWidgets)
// 	// widgetGroup.Get("/", handlers.Widget.GetWidgetByStatus)
// 	widgetGroup.Get("/:widget_id", handlers.Widget.GetWidget)
// 	widgetGroup.Get("/:widget_id/logs", handlers.Widget.GetLogs)
// 	widgetGroup.Put("/:widget_id", handlers.Widget.UpdateWidget)
// 	widgetGroup.Patch("/:widget_id/status", handlers.Widget.ChangeStatus)
// 	widgetGroup.Delete("/:widget_id", handlers.Widget.DeleteWidget)

// 	widgetGroup.Post("/", handlers.Widget.CreateWidget)
// 	widgetGroup.Post("/:widgetId/command", handlers.Command.SendCommand)

// 	// user
// 	userGroup := api.Group("/users")

// 	userGroup.Get("/", handlers.User.ListUsers)
// 	userGroup.Post("/", handlers.User.CreateUser)
// 	userGroup.Get("/:email", handlers.User.GetUser)
// 	userGroup.Delete("/:email", handlers.User.DeleteUser)
// 	// userGroup.Post("/:user_id/change-password", handlers.User.ChangePassword)
// 	// userGroup.Post("/:user_id/upload-profile", handlers.User.UploadProfile)
// }

func SetupRoutes(app *fiber.App, handlers *AppHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// api := app.Group("/api")

	// authen
	authGroup := app.Group("/auth")
	authGroup.Post("/login", handlers.User.Login)
	// authGroup.Post("/register", handlers.Auth.Register)

	adminOnly := middleware.RoleGuard(handlers.UserRepo, "ADMIN")

	
	authMid := middleware.FirebaseAuth(handlers.AuthClient,handlers.UserRepo)
	// api := app.Group("/api",authMid)
	api := app.Group("/api",authMid)

	secretGroup := api.Group("/secret", middleware.JWTMiddleware())
	secretGroup.Get("/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"secret_data": "นี่คือข้อมูลลับที่ได้รับการป้องกันด้วย JWT",
		})
	})

	// device
	deviceGroup := api.Group("/devices")
	deviceGroup.Get("/", handlers.Device.ListDevices)
	deviceGroup.Get("/:device_id", handlers.Device.GetDevice)
	deviceGroup.Put("/:device_id",adminOnly, handlers.Device.UpdateDevice)
	deviceGroup.Post("/:device_id/pair",adminOnly, handlers.Device.PairDevice)
	deviceGroup.Post("/:device_id/unpair",adminOnly, handlers.Device.UnpairDevice)

	// room
	roomGroup := api.Group("/rooms")
	roomGroup.Get("/", handlers.Room.ListRooms)
	roomGroup.Post("/",adminOnly, handlers.Room.CreateRoom)
	roomGroup.Get("/:room_id", handlers.Room.GetRoom)
	roomGroup.Put("/:room_id", handlers.Room.UpdateRoom)
	roomGroup.Delete("/:room_id",adminOnly, handlers.Room.DeleteRoom)
	roomGroup.Delete("/:room_id",adminOnly, handlers.Room.DeleteRoom)

	roomGroup.Post("/:room_id/devices",adminOnly, handlers.Room.AddDevice)
	roomGroup.Get("/:room_id/devices",adminOnly, handlers.Room.ListDevices)
	roomGroup.Post("/:room_id/widgets/order",adminOnly, handlers.Widget.ChangeOrder)

	roomGroup.Get("/:room_id/widgets", handlers.Widget.ListByRoom)

	// widget
	widgetGroup := api.Group("/widgets")

	widgetGroup.Get("/", handlers.Widget.ListWidgets)
	// widgetGroup.Get("/", handlers.Widget.GetWidgetByStatus)
	widgetGroup.Get("/:widget_id", handlers.Widget.GetWidget)
	widgetGroup.Get("/:widget_id/logs", handlers.Widget.GetLogs)
	widgetGroup.Put("/:widget_id",adminOnly, handlers.Widget.UpdateWidget)
	widgetGroup.Patch("/:widget_id/status",adminOnly, handlers.Widget.ChangeStatus)
	widgetGroup.Delete("/:widget_id",adminOnly, handlers.Widget.DeleteWidget)

	widgetGroup.Post("/", handlers.Widget.CreateWidget)
	widgetGroup.Post("/:widgetId/command", handlers.Command.SendCommand)

	// user
	userGroup := api.Group("/users")

	userGroup.Get("/", handlers.User.ListUsers)
	userGroup.Post("/",adminOnly, handlers.User.CreateUser)
	userGroup.Get("/:email", handlers.User.GetUser)
	userGroup.Delete("/:email",adminOnly, handlers.User.DeleteUser)
	// userGroup.Post("/:user_id/change-password", handlers.User.ChangePassword)
	// userGroup.Post("/:user_id/upload-profile", handlers.User.UploadProfile)
}