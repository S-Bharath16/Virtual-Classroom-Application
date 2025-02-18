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

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// deptID validation
	if student.DeptID != nil {
		var exists bool
		checkDeptQuery := `SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)`
		err = dbConn.QueryRow(checkDeptQuery, student.DeptID).Scan(&exists)
		if err != nil {
			log.Printf("Error checking department existence: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error checking if department exists",
			})
		}
		if !exists {
			log.Printf("Department ID %d does not exist", *student.DeptID)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid department ID. Please check and try again.",
			})
		}
	}

	// insert query
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