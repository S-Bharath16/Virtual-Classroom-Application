package modules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterDepartment(c *fiber.Ctx) error {
	var dept models.Department
	if err := c.BodyParser(&dept); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if dept.DepartmentName == "" {
		log.Println("Validation failed: department name is required")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Department name is required",
		})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Check if  dept  exists
	checkQuery := `
		SELECT department_id FROM deptData WHERE department_name = $1
	`
	var existingID int
	err = dbConn.QueryRow(checkQuery, dept.DepartmentName).Scan(&existingID)
	if err == nil {
		// If a row is found
		log.Printf("Department %s already exists with ID: %d", dept.DepartmentName, existingID)
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Department already exists",
		})
	}

	// Insert  new department
	insertQuery := `
		INSERT INTO deptData (deptName)
		VALUES ($1)
		RETURNING deptID
	`
	var newDeptID int
	err = dbConn.QueryRow(insertQuery, dept.DepartmentName).Scan(&newDeptID)
	if err != nil {
		log.Printf("Error inserting department: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create department",
		})
	}

	dept.DepartmentID = uint(newDeptID)
	log.Printf("Department created with ID: %d", newDeptID)

	return c.Status(http.StatusCreated).JSON(dept)
}
