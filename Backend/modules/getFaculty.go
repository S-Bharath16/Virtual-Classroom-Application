package modules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)


func GetAllFaculty(c *fiber.Ctx) error {
	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Select query for fetching
	query := `SELECT facultyID, facultyName, deptID FROM facultyData`
	rows, err := dbConn.Query(query)
	if err != nil {
		log.Printf("Error querying faculty data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query faculty data",
		})
	}
	defer rows.Close()

	facultyList := []models.Faculty{}
	for rows.Next() {
		var faculty models.Faculty
		if err := rows.Scan(&faculty.FacultyID, &faculty.FacultyName, &faculty.DeptID); err != nil {
			log.Printf("Error scanning faculty data: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse faculty data",
			})
		}
		facultyList = append(facultyList, faculty)
	}

	return c.Status(http.StatusOK).JSON(facultyList)
}
