package routes

import (
	"Backend/modules/Facultymodules"

	"github.com/gofiber/fiber/v2"
)

func RegisterFacultyRoutes(app *fiber.App) {
	api := app.Group("/api")
	faculty := api.Group("/faculty")
	faculty.Post("/createQuiz", Facultymodules.CreateQuiz)
	faculty.Put("/updateFaculty", Facultymodules.UpdateFaculty)
	faculty.Post("/recordAttendance", Facultymodules.RecordAttendance);
	faculty.Post("/createMeeting", Facultymodules.CreateMeeting);
	faculty.Post("/addAssignment", Facultymodules.AddAssignment)
}