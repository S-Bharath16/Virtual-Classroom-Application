package modules

import (
	"log"
	"net/http"
	"Backend/database"
	"Backend/models"
	"github.com/gofiber/fiber/v2"
)

func RegisterStudent(c *fiber.Ctx) error {
	var student models.Student
	if err := c.BodyParser(&student); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Data validation
	if student.RollNumber == "" || student.EmailID == "" || student.StudentName == "" ||
		student.StartYear == 0 || student.EndYear == 0 || student.Section == "" || student.Semester == 0 {
		log.Println("Validation failed: missing required fields")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Get the underlying sql.DB from GORM
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// SQL query using the exact column names
	query := `
		INSERT INTO studentData 
		(rollNumber, emailID, studentName, startYear, endYear, deptID, section, semester)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING studentID
	`

	var newID int
	err = dbConn.QueryRow(query,
		student.RollNumber,
		student.EmailID,
		student.StudentName,
		student.StartYear,
		student.EndYear,
		student.DeptID,
		student.Section,
		student.Semester,
	).Scan(&newID)
	if err != nil {
		log.Printf("[ERROR]: Error inserting student data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create student",
		})
	}

	return c.Status(http.StatusCreated).JSON(student)
}