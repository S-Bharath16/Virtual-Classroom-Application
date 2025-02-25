package Adminmodules

import (
	"database/sql"
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterCourse(c *fiber.Ctx) error {
	var course models.Course
	if err := c.BodyParser(&course); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	//  required fields
	if course.CourseCode == "" || course.CourseName == "" || course.CourseDeptID == 0 || course.CourseType == "" {
		log.Println("Validation failed: missing required fields")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "courseCode, courseName, courseDeptID and courseType are required",
		})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// checking if dept exists
	var deptExists bool
	checkDeptQuery := `SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)`
	err = dbConn.QueryRow(checkDeptQuery, course.CourseDeptID).Scan(&deptExists)
	if err != nil {
		log.Printf("Error checking department existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking department existence",
		})
	}
	if !deptExists {
		log.Printf("Department ID %d does not exist", course.CourseDeptID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid courseDeptID. Please check and try again.",
		})
	}

	// checking if admin exists
	if course.UpdatedBy != nil {
		var adminExists bool
		checkAdminQuery := `SELECT EXISTS (SELECT 1 FROM adminData WHERE adminID = $1)`
		err = dbConn.QueryRow(checkAdminQuery, course.UpdatedBy).Scan(&adminExists)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking admin existence: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error checking admin existence",
			})
		}
		if !adminExists {
			log.Printf("Admin ID %d does not exist", *course.UpdatedBy)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid updatedBy adminID. Please check and try again.",
			})
		}
	}

	// Insert
	query := `
		INSERT INTO courseData 
		(courseCode, courseName, courseDeptID, courseType, updatedBy)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING courseID
	`
	var newID int
	err = dbConn.QueryRow(query,
		course.CourseCode,
		course.CourseName,
		course.CourseDeptID,
		course.CourseType,
		course.UpdatedBy,
	).Scan(&newID)
	if err != nil {
		log.Printf("Error inserting course data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create course",
		})
	}

	course.CourseID = uint(newID)
	return c.Status(http.StatusCreated).JSON(course)
}
