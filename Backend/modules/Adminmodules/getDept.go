package Adminmodules

import (
	"log"
	"fmt"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)


func GetAllDepartments(c *fiber.Ctx) error {
	fmt.Println("[LOG]: Received Request From Load Tester :)");
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Select query
	query := `SELECT deptID, deptName FROM deptData`
	rows, err := dbConn.Query(query)
	if err != nil {
		log.Printf("Error querying departments: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query departments",
		})
	}
	defer rows.Close()

	departments := []models.Department{}
	for rows.Next() {
		var dept models.Department
		if err := rows.Scan(&dept.DepartmentID, &dept.DepartmentName); err != nil {
			log.Printf("Error scanning department: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse departments",
			})
		}
		departments = append(departments, dept)
	}

	return c.Status(http.StatusOK).JSON(departments)
}
