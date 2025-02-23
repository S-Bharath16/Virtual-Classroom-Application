package modules

import (
	"log"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func AssignFaculty(c *fiber.Ctx) error {
	var cf models.CourseFaculty
	if err := c.BodyParser(&cf); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation: required fields must be provided.
	if cf.CourseID == 0 || cf.FacultyID == 0 || cf.SectionID == 0 || cf.SemesterID == 0 || cf.DeptID == 0 {
		log.Println("Validation failed: missing required fields")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "courseID, facultyID, sectionID, semesterID and deptID are required",
		})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Validate if course exists.
	var courseExists bool
	courseQuery := `SELECT EXISTS (SELECT 1 FROM courseData WHERE courseID = $1)`
	if err := dbConn.QueryRow(courseQuery, cf.CourseID).Scan(&courseExists); err != nil {
		log.Printf("Error checking course existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking course existence",
		})
	}
	if !courseExists {
		log.Printf("Course ID %d does not exist", cf.CourseID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid courseID. Please check and try again.",
		})
	}

	// Validate if faculty exists.
	var facultyExists bool
	facultyQuery := `SELECT EXISTS (SELECT 1 FROM facultyData WHERE facultyID = $1)`
	if err := dbConn.QueryRow(facultyQuery, cf.FacultyID).Scan(&facultyExists); err != nil {
		log.Printf("Error checking faculty existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking faculty existence",
		})
	}
	if !facultyExists {
		log.Printf("Faculty ID %d does not exist", cf.FacultyID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid facultyID. Please check and try again.",
		})
	}

	// Validate if section exists.
	var sectionExists bool
	sectionQuery := `SELECT EXISTS (SELECT 1 FROM sectionData WHERE sectionID = $1)`
	if err := dbConn.QueryRow(sectionQuery, cf.SectionID).Scan(&sectionExists); err != nil {
		log.Printf("Error checking section existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking section existence",
		})
	}
	if !sectionExists {
		log.Printf("Section ID %d does not exist", cf.SectionID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sectionID. Please check and try again.",
		})
	}

	// Validate if semester exists.
	var semesterExists bool
	semesterQuery := `SELECT EXISTS (SELECT 1 FROM semesterData WHERE semesterID = $1)`
	if err := dbConn.QueryRow(semesterQuery, cf.SemesterID).Scan(&semesterExists); err != nil {
		log.Printf("Error checking semester existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking semester existence",
		})
	}
	if !semesterExists {
		log.Printf("Semester ID %d does not exist", cf.SemesterID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid semesterID. Please check and try again.",
		})
	}

	// Validate if deptID exists.
	var deptExists bool
	deptQuery := `SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)`
	if err := dbConn.QueryRow(deptQuery, cf.DeptID).Scan(&deptExists); err != nil {
		log.Printf("Error checking department existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking department existence",
		})
	}
	if !deptExists {
		log.Printf("Department ID %d does not exist", cf.DeptID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid deptID. Please check and try again.",
		})
	}

	// Validate basic rule:
	// "One section in one semester in one department can have only one faculty for a particular course.
	// Two faculties can't handle the same course for the same section in the same department.
	// (A faculty can take the same course across multiple sections, and can also take different courses across sections in different semesters/departments.)"
	var alreadyAssigned bool
	assignmentQuery := `
		SELECT EXISTS (
			SELECT 1 FROM courseFaculty 
			WHERE courseID = $1 
			  AND sectionID = $2 
			  AND semesterID = $3
			  AND deptID = $4
		)
	`
	if err := dbConn.QueryRow(assignmentQuery, cf.CourseID, cf.SectionID, cf.SemesterID, cf.DeptID).Scan(&alreadyAssigned); err != nil {
		log.Printf("Error checking existing assignment: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking existing assignment",
		})
	}
	if alreadyAssigned {
		log.Printf("Assignment already exists for courseID %d in sectionID %d, semesterID %d, and deptID %d", cf.CourseID, cf.SectionID, cf.SemesterID, cf.DeptID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "A faculty is already assigned for this course in the specified section, semester, and department.",
		})
	}

	// Validate adminID for createdBy and updatedBy remains unchanged.
	if cf.CreatedBy != nil {
		var adminExists bool
		adminQuery := `SELECT EXISTS (SELECT 1 FROM adminData WHERE adminID = $1)`
		if err := dbConn.QueryRow(adminQuery, cf.CreatedBy).Scan(&adminExists); err != nil {
			log.Printf("Error checking createdBy admin existence: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error checking createdBy admin existence",
			})
		}
		if !adminExists {
			log.Printf("Admin ID %d (createdBy) does not exist", *cf.CreatedBy)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid createdBy adminID. Please check and try again.",
			})
		}
	}

	if cf.UpdatedBy != nil {
		var adminExists bool
		adminQuery := `SELECT EXISTS (SELECT 1 FROM adminData WHERE adminID = $1)`
		if err := dbConn.QueryRow(adminQuery, cf.UpdatedBy).Scan(&adminExists); err != nil {
			log.Printf("Error checking updatedBy admin existence: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error checking updatedBy admin existence",
			})
		}
		if !adminExists {
			log.Printf("Admin ID %d (updatedBy) does not exist", *cf.UpdatedBy)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid updatedBy adminID. Please check and try again.",
			})
		}
	}

	// Insert the courseFaculty record (including the mandatory deptID).
	insertQuery := `
		INSERT INTO courseFaculty 
		(courseID, facultyID, sectionID, semesterID, deptID, createdBy, updatedBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING classroomID
	`
	var newID int
	err = dbConn.QueryRow(insertQuery,
		cf.CourseID,
		cf.FacultyID,
		cf.SectionID,
		cf.SemesterID,
		cf.DeptID,
		cf.CreatedBy,
		cf.UpdatedBy,
	).Scan(&newID)
	if err != nil {
		log.Printf("Error inserting courseFaculty data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to assign faculty",
		})
	}

	cf.ClassroomID = uint(newID)
	return c.Status(http.StatusCreated).JSON(cf)
}
