package Studentmodules

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// Request body struct
type StudentemailRequest struct {
	EmailID string `json:"studentEmail"`
}

func GetCourses(c *fiber.Ctx) error {
	// Parse request body
	var request StudentemailRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if request.EmailID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "emailID is required"})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Fetch sectionID, semesterID, and deptID
	var sectionID, semesterID, deptID int
	queryStudent := `
		SELECT sectionID, semesterID, deptID 
		FROM studentData 
		WHERE emailID = $1`
	err = dbConn.QueryRow(queryStudent, request.EmailID).Scan(&sectionID, &semesterID, &deptID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Student not found",
			})
		}
		log.Printf("Error fetching student data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch student data",
		})
	}

	// Query to fetch classroomID, courseID, courseCode, courseName, and assigned faculty details
	query := `
		SELECT cf.classroomID, cd.courseID, cd.courseCode, cd.courseName, f.facultyID, f.facultyName
		FROM courseFaculty cf
		JOIN courseData cd ON cf.courseID = cd.courseID
		JOIN facultyData f ON cf.facultyID = f.facultyID
		WHERE cf.sectionID = $1
		  AND cf.semesterID = $2
		  AND cf.deptID = $3
	`

	rows, err := dbConn.Query(query, sectionID, semesterID, deptID)
	if err != nil {
		log.Printf("Error querying courses: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query courses",
		})
	}
	defer rows.Close()

	courses := []models.StudentCourse{}
	for rows.Next() {
		var sc models.StudentCourse
		if err := rows.Scan(&sc.ClassroomID, &sc.CourseID, &sc.CourseCode, &sc.CourseName, &sc.FacultyID, &sc.FacultyName); err != nil {
			log.Printf("Error scanning row: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse course data",
			})
		}
		courses = append(courses, sc)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error iterating through course rows",
		})
	}

	// Return the course data
	return c.Status(http.StatusOK).JSON(courses)
}
