package routes

import (
	"Backend/modules"
	"Backend/modules/Auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	api := app.Group("/api")
	student := api.Group("/student")
	student.Get("/getQuiz", modules.GetQuiz);
	student.Get("/profile", modules.GetStudentProfile);
	student.Get("/getCourses", modules.GetCourses)
	student.Get("/auth/microsoft/url", Auth.HandleMicrosoftURL);
	student.Get("/auth/microsoft/callback", Auth.HandleMicrosoftCallback);
}
