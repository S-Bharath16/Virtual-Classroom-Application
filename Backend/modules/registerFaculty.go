package modules
import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
	"database/sql"
)

func RegisterFaculty(c *fiber.Ctx) error {
	var faculty models.Faculty
	if err := c.BodyParser(&faculty); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if faculty.EmailID == "" || faculty.FacultyName == "" {
		log.Println("Validation failed: emailID and facultyName are required")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Faculty name are required",
		})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Validation of deptID
	if faculty.DeptID != nil {
		var exists bool
		checkDeptQuery := `SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)`
		err = dbConn.QueryRow(checkDeptQuery, faculty.DeptID).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking department existence: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error checking department existence",
			})
		}
		if !exists {
			log.Printf("Department ID %d does not exist", *faculty.DeptID)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid department ID. Please check and try again.",
			})
		}
	}

	// Insert query
	query := `
		INSERT INTO facultyData (emailID, facultyName, deptID)
		VALUES ($1, $2, $3)
		RETURNING facultyID
	`
	var newID int
	err = dbConn.QueryRow(query,
		faculty.EmailID,
		faculty.FacultyName,
		faculty.DeptID,
	).Scan(&newID)
	if err != nil {
		log.Printf("[ERROR]: Error inserting faculty data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create faculty",
		})
	}

	faculty.FacultyID = uint(newID)

	return c.Status(http.StatusCreated).JSON(faculty)
}