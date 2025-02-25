package Adminmodules

import (
	"database/sql"
	"log"
	"net/http"

	"Backend/database"
	//"Backend/models"

	"github.com/gofiber/fiber/v2"
)

type RemoveStudentInput struct {
	StudentName string `json:"studentName"`
	EmailID     string `json:"emailID"`
}

func RemoveStudent(c *fiber.Ctx) error {
	var input RemoveStudentInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if input.StudentName == "" || input.EmailID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Both studentName and emailID are required",
		})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Check if  student exists.
	var studentID uint
	checkQuery := `
		SELECT studentID 
		FROM studentData 
		WHERE emailID = $1 AND studentName = $2
	`
	err = dbConn.QueryRow(checkQuery, input.EmailID, input.StudentName).Scan(&studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Student not found",
			})
		}
		log.Printf("Error checking student existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking student existence",
		})
	}

	// Student exists, so delete record.
	deleteQuery := `
		DELETE FROM studentData 
		WHERE studentID = $1
		RETURNING studentID
	`
	var deletedID uint
	err = dbConn.QueryRow(deleteQuery, studentID).Scan(&deletedID)
	if err != nil {
		log.Printf("Error deleting student: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove student",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":   "Student removed successfully",
		"studentID": deletedID,
	})
}
