package routes

import (
	"Backend/modules"

	"github.com/gofiber/fiber/v2"
)

func RegisterFacultyRoutes(app *fiber.App) {
	api := app.Group("/api")
	faculty := api.Group("/faculty")
	faculty.Post("/createQuiz", modules.CreateQuiz)
}
