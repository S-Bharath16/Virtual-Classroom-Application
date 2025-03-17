package Adminmodules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetCourses(c *fiber.Ctx) error {
	dbConn, err := database.GetDB().DB();
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	query := `
		SELECT * FROM courseData
	`;

	rows, err := dbConn.Query(query);
	if err != nil {
		log.Println("Error in Querying DB: ", err);
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to Fetch Courses from DB",
		})
	}
	defer rows.Close();

	courseData := []models.Course{};
	for rows.Next() {
		var currCourse models.Course
		if err := rows.Scan(
			&currCourse.CourseCode,
			&currCourse.CourseDeptID,
			&currCourse.CourseID,
			&currCourse.CourseName,
			&currCourse.CourseType,
			&currCourse.CreatedAt,
			&currCourse.UpdatedAt,
			&currCourse.UpdatedBy,
		); err != nil {
			log.Printf("Error scanning Course Data: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse Courses",
			});
		}
		courseData = append(courseData, currCourse);
	}
	return c.Status(http.StatusOK).JSON(courseData);
}