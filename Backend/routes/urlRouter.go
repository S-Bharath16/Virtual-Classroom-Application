package routes

import (
	"Backend/modules/Auth"
	"github.com/gofiber/fiber/v2"
)

func UrlRouter(app *fiber.App) {
	api := app.Group("/api");
	api.Get("/auth/google/url", Auth.HandleGoogleURL);
	api.Get("/auth/google/callback", Auth.HandleGoogleCallback);
}