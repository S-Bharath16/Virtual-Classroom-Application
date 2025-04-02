package Adminmodules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)


func GetAllStudents(c *fiber.Ctx) error {
	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Select query for fetching
	query := `SELECT studentID, rollNumber, emailID, studentName, startYear, endYear, deptID, sectionID, semesterID 
		FROM studentData`
	rows, err := dbConn.Query(query)
	if err != nil {
		log.Printf("Error querying student data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query student data",
		})
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		var student models.Student
		if err := rows.Scan(
			&student.StudentID, 
			&student.RollNumber, 
			&student.EmailID, 
			&student.StudentName, 
			&student.StartYear, 
			&student.EndYear, 
			&student.DeptID, 
			&student.SectionID, 
			&student.SemesterID,
		); err != nil {
			log.Printf("Error scanning student data: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse student data",
			})
		}
		students = append(students, student)
	}

	return c.Status(http.StatusOK).JSON(students)
}
