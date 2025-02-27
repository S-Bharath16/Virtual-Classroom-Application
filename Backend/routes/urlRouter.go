package routes

import (
	"Backend/modules/Auth/student"
	"Backend/modules/Auth/admin"
	"github.com/gofiber/fiber/v2"
)

func UrlRouter(app *fiber.App) {
	api := app.Group("/api");
	api.Get("/auth/google/url", student.HandleGoogleURL);
	api.Get("/studentAuth/google/callback", student.HandleGoogleCallback);
	api.Get("/adminAuth/google/callback", admin.HandleGoogleCallback);
}