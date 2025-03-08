package routes

import (
	"Backend/modules/Studentmodules"
	// "Backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterStudentRoutes(app *fiber.App) {
	api := app.Group("/api")
	student := api.Group("/student") 
	// student.Use(middleware.WebTokenValidator)

	student.Get("/getQuiz", Studentmodules.GetQuiz)
	student.Get("/profile", Studentmodules.GetStudentProfile)
	student.Get("/getCourses", Studentmodules.GetCourses)
	student.Get("/getallQuizzes", Studentmodules.GetAllQuizzes)
	student.Get("/attendance", Studentmodules.GetStudentAttendance);
}
