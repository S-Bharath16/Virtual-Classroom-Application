package Studentmodules

import (
	"log"
	"net/http"
	"strconv"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)


func GetCourses(c *fiber.Ctx) error {
	// Retrieve, validate query parameters.
	sectionIDParam := c.Query("sectionID")
	semesterIDParam := c.Query("semesterID")
	deptIDParam := c.Query("deptID")

	if sectionIDParam == "" || semesterIDParam == "" || deptIDParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "sectionID, semesterID, and deptID query parameters are required",
		})
	}

	sectionID, err := strconv.Atoi(sectionIDParam)
	if err != nil || sectionID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sectionID parameter",
		})
	}

	semesterID, err := strconv.Atoi(semesterIDParam)
	if err != nil || semesterID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid semesterID parameter",
		})
	}

	deptID, err := strconv.Atoi(deptIDParam)
	if err != nil || deptID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid deptID parameter",
		})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	//  join query:  fetch the classroomID, courseID and courseName from courseData,
	// along with the assigned faculty details from facultyData.
	query := `
		SELECT cf.classroomID, cd.courseID, cd.courseName, f.facultyID, f.facultyName
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
		if err := rows.Scan(&sc.ClassroomID, &sc.CourseID, &sc.CourseName, &sc.FacultyID, &sc.FacultyName); err != nil {
			log.Printf("Error scanning row: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse course data",
			})
		}
		courses = append(courses, sc)
	}

	return c.Status(http.StatusOK).JSON(courses)
}
