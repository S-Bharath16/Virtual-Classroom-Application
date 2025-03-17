package routes

import (
	"Backend/modules/Adminmodules"
	"Backend/modules/mailer"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudent(app *fiber.App) {
	api := app.Group("/api")
	admin := api.Group("/admin")
	admin.Post("/registerStudent", Adminmodules.RegisterStudent);
	admin.Post("/registerDepartment", Adminmodules.RegisterDepartment);
	admin.Post("/registerFaculty", Adminmodules.RegisterFaculty);
	admin.Post("/registerCourse", Adminmodules.RegisterCourse);
	admin.Post("/assignFaculty", Adminmodules.AssignFaculty);
	admin.Get("/getAllDepartments", Adminmodules.GetAllDepartments);
	admin.Get("/getAllFaculty", Adminmodules.GetAllFaculty);
	admin.Post("/announcementMail", mailer.SendAnnouncementMail);
	admin.Put("/updateCourse", Adminmodules.UpdateCourse);
	admin.Delete("/removestudent", Adminmodules.RemoveStudent);
	admin.Delete("/removeCourse", Adminmodules.RemoveCourse);
}