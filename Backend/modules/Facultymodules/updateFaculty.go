package Facultymodules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func UpdateFaculty(c *fiber.Ctx) error {
	var faculty models.Faculty
	if err := c.BodyParser(&faculty); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if faculty.FacultyID == 0 || faculty.EmailID == "" || faculty.FacultyName == "" || faculty.DeptID == nil || *faculty.DeptID == 0 {
		log.Println("Validation failed: missing required fields")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "facultyID, emailID, facultyName and deptID are required",
		})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Validate facultyID .
	var facultyExists bool
	facultyCheckQuery := `SELECT EXISTS (SELECT 1 FROM facultyData WHERE facultyID = $1)`
	if err := dbConn.QueryRow(facultyCheckQuery, faculty.FacultyID).Scan(&facultyExists); err != nil {
		log.Printf("Error checking faculty existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking faculty existence",
		})
	}
	if !facultyExists {
		log.Printf("Faculty ID %d does not exist", faculty.FacultyID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Faculty with provided facultyID does not exist",
		})
	}

	// Validate deptID 
	var deptExists bool
	deptQuery := `SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)`
	if err := dbConn.QueryRow(deptQuery, faculty.DeptID).Scan(&deptExists); err != nil {
		log.Printf("Error checking department existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking department existence",
		})
	}
	if !deptExists {
		log.Printf("Department ID %d does not exist", *faculty.DeptID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid deptID. Please check and try again.",
		})
	}

	// Update faculty record query.
	updateQuery := `
		UPDATE facultyData
		SET emailID = $1,
		    facultyName = $2,
		    deptID = $3,
		    updatedAt = CURRENT_TIMESTAMP
		WHERE facultyID = $4
		RETURNING facultyID, emailID, facultyName, deptID, createdAt, updatedAt
	`
	var updatedFaculty models.Faculty
	err = dbConn.QueryRow(updateQuery,
		faculty.EmailID,
		faculty.FacultyName,
		faculty.DeptID,
		faculty.FacultyID,
	).Scan(
		&updatedFaculty.FacultyID,
		&updatedFaculty.EmailID,
		&updatedFaculty.FacultyName,
		&updatedFaculty.DeptID,
		&updatedFaculty.CreatedAt,
		&updatedFaculty.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error updating faculty: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update faculty",
		})
	}

	return c.Status(http.StatusOK).JSON(updatedFaculty)
}
