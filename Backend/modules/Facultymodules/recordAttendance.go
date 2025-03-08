package Facultymodules

import (
	"Backend/database"
	"Backend/models"
	"database/sql"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RecordAttendance(c *fiber.Ctx) error {
	// Get DB connection
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	var attendance models.Attendance;

	if err := c.BodyParser(&attendance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var meetingID int
	err = dbConn.QueryRow("SELECT meetingID FROM meetingData WHERE meetingID = $1", attendance.MeetingID).Scan(&meetingID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid meetingID"})
	} else if err != nil {
		log.Printf("Error querying meetingID: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var studentID int
	err = dbConn.QueryRow("SELECT studentID FROM studentData WHERE studentID = $1", attendance.StudentID).Scan(&studentID)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid studentID"})
	} else if err != nil {
		log.Printf("Error querying studentID: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	var existingIsPresent string
	err = dbConn.QueryRow("SELECT isPresent FROM attendanceData WHERE meetingID = $1 AND studentID = $2", attendance.MeetingID, attendance.StudentID).Scan(&existingIsPresent)
	if err == nil {
		// Record exists, update it
		_, err = dbConn.Exec("UPDATE attendanceData SET isPresent = $1 WHERE meetingID = $2 AND studentID = $3", attendance.IsPresent, attendance.MeetingID, attendance.StudentID)
		if err != nil {
			log.Printf("Error updating attendance record: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update attendance"})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Attendance updated successfully"})
	} else if err != sql.ErrNoRows {
		log.Printf("Error checking existing attendance record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	_, err = dbConn.Exec("INSERT INTO attendanceData (meetingID, studentID, isPresent) VALUES ($1, $2, $3)", 
		attendance.MeetingID, attendance.StudentID, attendance.IsPresent)
	if err != nil {
		log.Printf("Error inserting attendance record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record attendance"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Attendance recorded successfully"})
}