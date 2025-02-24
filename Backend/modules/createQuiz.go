package modules

import (
	"log"
	"net/http"
	"time"

	"Backend/database"
	"Backend/models"
	"Backend/modules/mailer"

	"github.com/gofiber/fiber/v2"
)


func CreateQuiz(c *fiber.Ctx) error {
	var quiz models.Quiz
	if err := c.BodyParser(&quiz); err != nil {
		log.Printf("Error parsing quiz data: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validatation
	if quiz.ClassroomID == 0 || quiz.QuizName == "" || quiz.StartTime.IsZero() || quiz.EndTime.IsZero() || quiz.QuizDuration == 0 || quiz.CreatedBy == 0 {
		log.Println("Validation failed: missing required fields")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "classroomID, quizName, startTime, endTime, quizDuration and createdBy are required",
		})
	}

	
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Validation classroomID 
	var classroomExists bool
	checkClassroomQuery := `SELECT EXISTS (SELECT 1 FROM courseFaculty WHERE classroomID = $1)`
	err = dbConn.QueryRow(checkClassroomQuery, quiz.ClassroomID).Scan(&classroomExists)
	if err != nil {
		log.Printf("Error checking classroom existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking classroom existence",
		})
	}
	if !classroomExists {
		log.Printf("ClassroomID %d does not exist", quiz.ClassroomID)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid classroomID. Please check and try again.",
		})
	}

	// Validation facultyData table
	var facultyExists bool
	checkFacultyQuery := `SELECT EXISTS (SELECT 1 FROM facultyData WHERE facultyID = $1)`
	err = dbConn.QueryRow(checkFacultyQuery, quiz.CreatedBy).Scan(&facultyExists)
	if err != nil {
		log.Printf("Error checking faculty existence: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking faculty existence",
		})
	}
	if !facultyExists {
		log.Printf("Faculty ID %d does not exist", quiz.CreatedBy)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid createdBy facultyID. Please check and try again.",
		})
	}

	//  default values 
	if quiz.IsOpenForAll == "" {
		quiz.IsOpenForAll = "0"
	}
	// CreatedAt and UpdatedAt are set
	quiz.CreatedAt = time.Now()
	quiz.UpdatedAt = time.Now()

	// Insert query
	insertQuery := `
		INSERT INTO quizData 
		(classroomID, quizName, quizDescription, quizData, isOpenForAll, startTime, endTime, quizDuration, createdBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING quizID
	`

	var newQuizID int
	err = dbConn.QueryRow(insertQuery,
		quiz.ClassroomID,
		quiz.QuizName,
		quiz.QuizDescription,
		quiz.QuizData,
		quiz.IsOpenForAll,
		quiz.StartTime,
		quiz.EndTime,
		quiz.QuizDuration,
		quiz.CreatedBy,
	).Scan(&newQuizID)
	if err != nil {
		log.Printf("Error inserting quiz data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create quiz",
		})
	}

	quiz.QuizID = uint(newQuizID)
	mailer.QuizCreationMail("Quiz Creation Alert Mail", "", int(quiz.ClassroomID));
	return c.Status(http.StatusCreated).JSON(quiz)
}
