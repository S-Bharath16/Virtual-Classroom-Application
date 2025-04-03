package Facultymodules

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"Backend/database"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

type QuizRequest struct {
	ClassroomID     int       `json:"classroomID"`
	QuizName        string    `json:"quizName"`
	QuizDescription string    `json:"quizDescription"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	QuizDuration    int       `json:"quizDuration"`
	CreatedBy       int       `json:"createdBy"`
	IsOpenForAll    bool      `json:"isOpenForAll"`
	AccessToken     string    `json:"accessToken"`
}

func CreateQuiz(c *fiber.Ctx) error {
	var request QuizRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if request.ClassroomID == 0 || request.QuizName == "" || request.StartTime.IsZero() || request.EndTime.IsZero() || request.QuizDuration == 0 || request.CreatedBy == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	// Create Google Form
	formURL, err := createGoogleForm(request.QuizName, request.QuizDescription, request.AccessToken)
	if err != nil {
		log.Printf("Failed to create Google Form: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create Google Form",
		})
	}

	// DB insert
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("DB connection error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database connection failed",
		})
	}

	query := `
		INSERT INTO quizData (classroomID, quizName, quizDescription, quizData, isOpenForAll, startTime, endTime, quizDuration, createdBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING quizID
	`

	var quizID int
	err = dbConn.QueryRow(query,
		request.ClassroomID,
		request.QuizName,
		request.QuizDescription,
		formURL,
		fmt.Sprintf("%v", request.IsOpenForAll),
		request.StartTime,
		request.EndTime,
		request.QuizDuration,
		request.CreatedBy,
	).Scan(&quizID)

	if err != nil {
		log.Printf("Failed to insert quiz: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database insert failed",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":  "Quiz created successfully",
		"quizID":   quizID,
		"formUrl":  formURL,
	})
}

func createGoogleForm(title, description, accessToken string) (string, error) {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)

	service, err := forms.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", fmt.Errorf("error creating forms service: %v", err)
	}

	form := &forms.Form{
		Info: &forms.Info{
			Title:       title,
			Description: description,
		},
	}

	res, err := service.Forms.Create(form).Do()
	if err != nil {
		return "", fmt.Errorf("error creating form: %v", err)
	}

	formURL := fmt.Sprintf("https://docs.google.com/forms/d/%s/edit", res.FormId)
	return formURL, nil
}
