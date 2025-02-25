package Adminmodules

import (
	"database/sql"
	"log"
	"net/http"

	"Backend/database"
	//"Backend/models"

	"github.com/gofiber/fiber/v2"
)

type RemoveCourseInput struct {
	CourseCode string `json:"courseCode"`
}

func RemoveCourse(c *fiber.Ctx) error {
	var input RemoveCourseInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate courseCode 
	if input.CourseCode == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "courseCode is required",
		})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Check if course exists.
	var courseID uint
	checkQuery := `SELECT courseID FROM courseData WHERE courseCode = $1`
	err = dbConn.QueryRow(checkQuery, input.CourseCode).Scan(&courseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Course not found",
			})
		}
		log.Printf("Error checking course existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking course existence",
		})
	}

	// Course exists, so delete record.
	deleteQuery := `DELETE FROM courseData WHERE courseCode = $1 RETURNING courseID`
	var deletedID uint
	err = dbConn.QueryRow(deleteQuery, input.CourseCode).Scan(&deletedID)
	if err != nil {
		log.Printf("Error deleting course: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove course",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":  "Course removed successfully",
		"courseID": deletedID,
	})
}
