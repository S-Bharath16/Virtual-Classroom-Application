package Studentmodules

import (
	"database/sql"
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// Rule: Fetches the student profile by joining studentData with deptData, sectionData, and semesterData, based on the emailID sent in the request body.
func GetStudentProfile(c *fiber.Ctx) error {
	
	var request struct {
		EmailID string `json:"emailID"`
	}

	// Parse 
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate emailID
	if request.EmailID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "emailID is required"})
	}

	// Get database connection
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get database connection"})
	}

	// Join query
	query := `
		SELECT 
			s.rollNumber, 
			s.emailID, 
			s.studentName, 
			s.startYear, 
			s.endYear, 
			d.deptName, 
			sec.sectionName, 
			sem.semesterNumber
		FROM studentData s
		LEFT JOIN deptData d ON s.deptID = d.deptID
		LEFT JOIN sectionData sec ON s.sectionID = sec.sectionID
		LEFT JOIN semesterData sem ON s.semesterID = sem.semesterID
		WHERE s.emailID = $1
	`

	var profile models.StudentProfile
	err = dbConn.QueryRow(query, request.EmailID).Scan(
		&profile.RollNumber,
		&profile.EmailID,
		&profile.StudentName,
		&profile.StartYear,
		&profile.EndYear,
		&profile.DeptName,
		&profile.SectionName,
		&profile.SemesterNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
		}
		log.Printf("Error querying student profile: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving student profile"})
	}

	// Return the student profile
	return c.Status(http.StatusOK).JSON(profile)
}
