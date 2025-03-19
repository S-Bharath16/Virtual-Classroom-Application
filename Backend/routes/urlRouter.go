package routes

import (
	"Backend/modules/Auth/student"
	"Backend/modules/Auth/admin"
	"Backend/modules/Auth/faculty"
	"github.com/gofiber/fiber/v2"
)

func UrlRouter(app *fiber.App) {
	api := app.Group("/api") // Root API group

	// Define studentAuth routes separately
	studentAuthRouter := api.Group("/studentAuth") 
	adminRouter := api.Group("/Aauth")
	facultyRouter := api.Group("/Fauth")

	api.Get("/auth/google/url", student.HandleGoogleURL)
	studentAuthRouter.Get("/google/callback", student.HandleGoogleCallback)
	adminRouter.Get("/google/callback", admin.HandleGoogleCallback)
	facultyRouter.Get("/google/callback", faculty.HandleGoogleCallback);
}