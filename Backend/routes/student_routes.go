package routes

import (
	"Backend/modules/Studentmodules"
	"Backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	api := app.Group("/api") // Root API group

	// Define student routes **without middleware first**
	student := api.Group("/student") 

	// Apply middleware **only to exact student routes**
	student.Use(middleware.WebTokenValidator)

	student.Get("/getQuiz", Studentmodules.GetQuiz)
	student.Get("/profile", Studentmodules.GetStudentProfile)
	student.Get("/getCourses", Studentmodules.GetCourses)
	student.Get("/getallQuizzes", Studentmodules.GetAllQuizzes)
}
