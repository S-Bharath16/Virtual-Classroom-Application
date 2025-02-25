package Adminmodules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func UpdateCourse(c *fiber.Ctx) error {
	var course models.Course
	if err := c.BodyParser(&course); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if course.CourseID == 0 || course.CourseCode == "" || course.CourseName == "" || course.CourseDeptID == 0 || course.CourseType == "" {
		log.Println("Validation failed: missing required fields")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "courseID, courseCode, courseName, courseDeptID and courseType are required",
		})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Validate  courseID 
	var courseExists bool
	courseCheckQuery := `SELECT EXISTS (SELECT 1 FROM courseData WHERE courseID = $1)`
	if err := dbConn.QueryRow(courseCheckQuery, course.CourseID).Scan(&courseExists); err != nil {
		log.Printf("Error checking course existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking course existence",
		})
	}
	if !courseExists {
		log.Printf("Course ID %d does not exist", course.CourseID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid courseID. Please check and try again.",
		})
	}

	// Validate  courseDeptID 
	var deptExists bool
	deptQuery := `SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)`
	if err := dbConn.QueryRow(deptQuery, course.CourseDeptID).Scan(&deptExists); err != nil {
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

	//  validate admin 
	if course.UpdatedBy != nil {
		var adminExists bool
		adminQuery := `SELECT EXISTS (SELECT 1 FROM adminData WHERE adminID = $1)`
		if err := dbConn.QueryRow(adminQuery, course.UpdatedBy).Scan(&adminExists); err != nil {
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

	// Update record
	updateQuery := `
		UPDATE courseData
		SET courseCode = $1,
		    courseName = $2,
		    courseDeptID = $3,
		    courseType = $4,
		    updatedBy = $5,
		    updatedAt = CURRENT_TIMESTAMP
		WHERE courseID = $6
		RETURNING courseID, courseCode, courseName, courseDeptID, createdAt, updatedAt, courseType, updatedBy
	`
	var updatedCourse models.Course
	err = dbConn.QueryRow(updateQuery,
		course.CourseCode,
		course.CourseName,
		course.CourseDeptID,
		course.CourseType,
		course.UpdatedBy,
		course.CourseID,
	).Scan(
		&updatedCourse.CourseID,
		&updatedCourse.CourseCode,
		&updatedCourse.CourseName,
		&updatedCourse.CourseDeptID,
		&updatedCourse.CreatedAt,
		&updatedCourse.UpdatedAt,
		&updatedCourse.CourseType,
		&updatedCourse.UpdatedBy,
	)
	if err != nil {
		log.Printf("Error updating course: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update course",
		})
	}

	return c.Status(http.StatusOK).JSON(updatedCourse)
}
