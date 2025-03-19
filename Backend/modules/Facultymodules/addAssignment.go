package Facultymodules

import (
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"fmt"
	"io"
	"net/http"

	"Backend/database"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func AddAssignment(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Error parsing form data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

    // Function to read file field as text
    readFileFieldAsText := func(fieldName string) (string, error) {
        files, ok := form.File[fieldName]
        if !ok || len(files) == 0 {
            return "", fmt.Errorf("field %s not found", fieldName)
        }
        
        file, err := files[0].Open()
        if err != nil {
            return "", err
        }
        defer file.Close()
        
        content, err := io.ReadAll(file)
        if err != nil {
            return "", err
        }
        
        return string(content), nil
    }

    // Extract fields from file parts
    classroomID, _ := readFileFieldAsText("classroomID")
    assignmentName, _ := readFileFieldAsText("assignmentName")
    assignmentDescription, _ := readFileFieldAsText("assignmentDescription")
    assignmentText, _ := readFileFieldAsText("assignmentText")
    createdBy, _ := readFileFieldAsText("createdBy")
    startTime, _ := readFileFieldAsText("startTime")
    endTime, _ := readFileFieldAsText("endTime")

	if classroomID == "" || assignmentName == "" || assignmentDescription == "" || assignmentText == "" || startTime == "" || endTime == "" || createdBy == "" {
		log.Println("Validation failed: missing required fields")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields (classroomID, assignmentName, assignmentDescription, assignmentText, startTime, endTime, submissionLimit, createdBy) are required",
		})
	}

	// Parse time fields
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid startTime format",
		})
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid endTime format",
		})
	}

	if start.After(end) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "startTime must be before endTime",
		})
	}

	// 2. Handle file uploads (REQUIRED)
	var imagePath, documentPath string

	// Handle Image Upload (REQUIRED)
	imageFile, err := getFormFile(form, "image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image file is required",
		})
	}
	imagePath, err = saveFile(c, imageFile, "images")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image file",
		})
	}

	// Handle Document Upload (REQUIRED)
	docFile, err := getFormFile(form, "document")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Document file is required",
		})
	}
	documentPath, err = saveFile(c, docFile, "documents")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save document file",
		})
	}

	assignment := models.Assignment{
		ClassroomID:          parseUint(classroomID),
		AssignmentName:       assignmentName,
		AssignmentDescription: assignmentDescription,
		AssignmentText:       assignmentText,
		ImagePath:            &imagePath,
		DocumentPath:         &documentPath,
		StartTime:            start,
		EndTime:              end,
		CreatedBy:            parseUint(createdBy),
	}

	dbConn, err := database.GetDB().DB()
    if err != nil {
        log.Printf("Error getting DB connection: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to get database connection",
        })
    }

	query := `
		INSERT INTO assignmentData (classroomID, assignmentName, assignmentDescription, assignmentText, imagePath, documentPath, startTime, endTime, createdBy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING assignmentID
	`

    var assignmentID int
    err = dbConn.QueryRow(query, assignment.ClassroomID, assignment.AssignmentName, assignment.AssignmentDescription, assignment.AssignmentText, assignment.ImagePath, assignment.DocumentPath, assignment.StartTime, assignment.EndTime, assignment.CreatedBy).Scan(&assignmentID)
    if err != nil {
        log.Printf("Error inserting meeting data: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to insert meeting data",
        })
    }

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Assignment created successfully",
		"assignmentID": assignmentID,
	})
}

// ✅ Extract Form File
func getFormFile(form *multipart.Form, key string) (*multipart.FileHeader, error) {
	if fileHeaders, ok := form.File[key]; ok && len(fileHeaders) > 0 {
		return fileHeaders[0], nil
	}
	return nil, fmt.Errorf("%s file not provided", key)
}

// ✅ Save File to Local Storage
func saveFile(c *fiber.Ctx, file *multipart.FileHeader, fileType string) (string, error) {
	dir := filepath.Join("uploads", fileType)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// Generate File Path
	filePath := filepath.Join(dir, filepath.Base(file.Filename))
	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return "/" + filePath, nil
}

// ✅ Convert String to Uint
func parseUint(value string) uint {
	var result uint
	if value != "" {
		fmt.Sscanf(value, "%d", &result)
	}
	return result
}