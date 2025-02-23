package routes

import (
	"Backend/modules"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	api := app.Group("/api")
	student := api.Group("/student")
	student.Get("/getQuiz", modules.GetQuiz);
	student.Get("/profile", modules.GetStudentProfile);
}
