package Studentmodules

import (
	"log"
	"net/http"
	"Backend/database"
	"github.com/gofiber/fiber/v2"
)

func GetMeetingLink(c *fiber.Ctx) error {
	var request struct {
		ClassroomID int `json:"classroomID"`
		StudentEmail	  string 	`json:"emailID"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	query := `SELECT meetingLink FROM meetingData WHERE classroomID = $1 ORDER BY startTime DESC LIMIT 1`
	var meetingLink string

	err = dbConn.QueryRow(query, request.ClassroomID).Scan(&meetingLink)
	if err != nil {
		log.Printf("Error fetching meeting link: %v", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "No meeting found for the given classroom ID",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Metting Link Returned",
		"meetingLink": meetingLink,
	});
}