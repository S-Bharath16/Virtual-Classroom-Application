package handlers

import (
	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func PutItem(c *fiber.Ctx) error {
	item := new(models.Item)
	if err := c.BodyParser(item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	result := database.DB.Create(&item)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save item to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(item)
}
