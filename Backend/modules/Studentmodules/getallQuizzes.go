package Studentmodules

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)


type StudentEmailRequest struct {
	EmailID string `json:"emailID"`
}

func GetAllQuizzes(c *fiber.Ctx) error {
	
	var request StudentEmailRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if request.EmailID == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "emailID is required"})
	}

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection error"})
	}

	// Retrieve student details
	var sectionID, semesterID, deptID int
	err = dbConn.QueryRow("SELECT sectionID, semesterID, deptID FROM studentData WHERE emailID = $1", request.EmailID).Scan(&sectionID, &semesterID, &deptID)
	if err != nil {
		log.Printf("Student not found for email: %v", request.EmailID)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
	}

	// Get current time
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

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		if err := rows.Scan(
			&quiz.QuizID, &quiz.ClassroomID, &quiz.QuizName, &quiz.QuizDescription, &quiz.QuizData,
			&quiz.IsOpenForAll, &quiz.StartTime, &quiz.EndTime, &quiz.QuizDuration,
			&quiz.CreatedAt, &quiz.UpdatedAt, &quiz.CreatedBy,
		); err != nil {
			log.Printf("Error scanning quiz row: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning quiz data"})
		}
		quizzes = append(quizzes, quiz)
	}

	return c.Status(http.StatusOK).JSON(quizzes)
}
