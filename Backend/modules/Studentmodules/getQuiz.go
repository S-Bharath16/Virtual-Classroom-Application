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


func GetQuiz(c *fiber.Ctx) error {
	// Get classroomID 
	classroomIDParam := c.Query("classroomID")
	if classroomIDParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "classroomID query parameter is required",
		})
	}

	classroomID, err := strconv.Atoi(classroomIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid classroomID parameter",
		})
	}

	// Get  time 
	currentTime := time.Now()

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Query
	query := `
		SELECT quizID, classroomID, quizName, quizDescription, quizData, isOpenForAll, startTime, endTime, quizDuration, createdAt, updatedAt, createdBy
		FROM quizData
		WHERE classroomID = $1 
		  AND isOpenForAll = '1'
		  AND startTime <= $2
		  AND endTime >= $2
	`
	rows, err := dbConn.Query(query, classroomID, currentTime)
	if err != nil {
		log.Printf("Error querying quiz data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query quiz data",
		})
	}
	defer rows.Close()

	quizzes := []models.Quiz{}
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
			log.Printf("Error scanning quiz: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse quiz data",
			})
		}
		quizzes = append(quizzes, quiz)
	}

	return c.Status(http.StatusOK).JSON(quizzes)
}
