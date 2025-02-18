package routes

import (
	"Backend/modules"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudent(app *fiber.App) {
	api := app.Group("/api")
	admin := api.Group("/admin")
	admin.Post("/registerStudent", modules.RegisterStudent);
}