package Studentmodules

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// Rule : retreive all active and upcoming quizzes for a stuent's respective section , semester and department.
func GetAllQuizzes(c *fiber.Ctx) error {
	
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sectionID"})
	}
	semesterID, err := strconv.Atoi(semesterIDParam)
	if err != nil || semesterID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid semesterID"})
	}
	deptID, err := strconv.Atoi(deptIDParam)
	if err != nil || deptID <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deptID"})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection error"})
	}

	// Validate section exists.
	var exists bool
	err = dbConn.QueryRow("SELECT EXISTS (SELECT 1 FROM sectionData WHERE sectionID = $1)", sectionID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sectionID"})
	}

	// Validate semester exists.
	err = dbConn.QueryRow("SELECT EXISTS (SELECT 1 FROM semesterData WHERE semesterID = $1)", semesterID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid semesterID"})
	}

	// Validate department exists.
	err = dbConn.QueryRow("SELECT EXISTS (SELECT 1 FROM deptData WHERE deptID = $1)", deptID).Scan(&exists)
	if err != nil || !exists {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deptID"})
	}

	// Get current time.
	currentTime := time.Now()

	
	query := `
		SELECT 
			q.quizID, q.classroomID, q.quizName, q.quizDescription, q.quizData, q.isOpenForAll,
			q.startTime, q.endTime, q.quizDuration, q.createdAt, q.updatedAt, q.createdBy
		FROM courseFaculty cf
		JOIN quizData q ON cf.classroomID = q.classroomID
		WHERE cf.sectionID = $1
		  AND cf.semesterID = $2
		  AND cf.deptID = $3
		  AND q.endTime >= $4
	`
	rows, err := dbConn.Query(query, sectionID, semesterID, deptID, currentTime)
	if err != nil {
		log.Printf("Error querying quizzes: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying quizzes"})
	}
	defer rows.Close()

	// Build  list of quizzes.
	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		err := rows.Scan(
			&quiz.QuizID,
			&quiz.ClassroomID,
			&quiz.QuizName,
			&quiz.QuizDescription,
			&quiz.QuizData,
			&quiz.IsOpenForAll,
			&quiz.StartTime,
			&quiz.EndTime,
			&quiz.QuizDuration,
			&quiz.CreatedAt,
			&quiz.UpdatedAt,
			&quiz.CreatedBy,
		)
		if err != nil {
			log.Printf("Error scanning quiz row: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning quiz data"})
		}
		quizzes = append(quizzes, quiz)
	}

	return c.Status(http.StatusOK).JSON(quizzes)
}