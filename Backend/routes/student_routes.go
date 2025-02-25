package routes

import (
	"Backend/modules/Studentmodules"
	"Backend/modules/Auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	api := app.Group("/api")
	student := api.Group("/student")
	student.Get("/getQuiz", Studentmodules.GetQuiz);
	student.Get("/profile", Studentmodules.GetStudentProfile);
	student.Get("/getCourses", Studentmodules.GetCourses)
	student.Get("/auth/microsoft/url", Auth.HandleMicrosoftURL);
	student.Get("/auth/microsoft/callback", Auth.HandleMicrosoftCallback);
}
