package routes

import (
	"Backend/modules"
	"Backend/modules/mailer"

	"github.com/gofiber/fiber/v2"
)

func RegisterStudent(app *fiber.App) {
	api := app.Group("/api")
	admin := api.Group("/admin")
	admin.Post("/registerStudent", modules.RegisterStudent);
	admin.Post("/registerDepartment", modules.RegisterDepartment);
	admin.Post("/registerFaculty", modules.RegisterFaculty);
	admin.Post("/registerCourse", modules.RegisterCourse);
	admin.Post("/assignFaculty", modules.AssignFaculty);
	admin.Get("/getAllDepartments", modules.GetAllDepartments);
	admin.Get("/getAllFaculty", modules.GetAllFaculty);
	admin.Post("/announcementMail", mailer.SendAnnouncementMail);
	admin.Put("/updateCourse", modules.UpdateCourse);
	admin.Delete("/removestudent", modules.RemoveStudent);
	admin.Delete("/removeCourse", modules.RemoveCourse);
}
