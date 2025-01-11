package routes

import (
	"Backend/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterItemRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Put("/item", handlers.PutItem)
}
