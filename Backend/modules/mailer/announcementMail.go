package mailer

import (
	"fmt"
	"log"
	"net/http"
	"Backend/database"
	"Backend/utilities/mailer"
	"github.com/gofiber/fiber/v2"
)

func SendAnnouncementMail(c *fiber.Ctx) error {

	var request struct {
		Subject  string   `json:"mailSubject"`
		Body     string   `json:"mailBody"`
		DeptIDs  []int    `json:"deptIDs"`
	}

	if err := c.BodyParser(&request); err != nil {
		log.Println("Failed to parse request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get database connection
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	// Query to fetch email IDs of students belonging to the specified deptIDs
	query := `
		SELECT s.emailID FROM studentData s
		WHERE s.deptID = ANY($1)
	`

	rows, err := dbConn.Query(query, request.DeptIDs)
	if err != nil {
		log.Println("Failed to retrieve student emails:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch student emails",
		})
	}
	defer rows.Close()

	var emailIDs []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Println("Error scanning email row:", err)
			continue
		}
		emailIDs = append(emailIDs, email)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over email rows:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing student emails",
		})
	}

	// Call SendMail function
	if err := mailer.SendMail(emailIDs, request.Subject, request.Body, "announcement"); err != nil {
		log.Println("Failed to send emails:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send emails",
		})
	}

	fmt.Println("[LOG]: Announcement Mail Sent Succussfully");

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Emails sent successfully",
	})
}