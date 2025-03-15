package Studentmodules

import (
	"Backend/models"
	"database/sql"
	"log"
	"net/http"
	"Backend/database"

	"github.com/gofiber/fiber/v2"
)

func GetStudentAttendance(c *fiber.Ctx) error {
	// Get DB connection
	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("Error getting DB connection: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get database connection",
		})
	}

	var requestData struct {
		StudentEmail string `json:"emailID"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Fetch Student ID
	var studentID int
	query := "SELECT studentID FROM studentData WHERE emailID = $1"
	if err := dbConn.QueryRow(query, requestData.StudentEmail).Scan(&studentID); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Student not found"})
		}
		log.Printf("Error fetching studentID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Fetch attendance records with total classes and numPresent
	attendanceQuery := `
		SELECT c.courseID, c.courseName, cf.classroomID,
			COUNT(a.meetingID) AS totalClasses,
			SUM(CASE WHEN a.isPresent = '1' THEN 1 ELSE 0 END) AS numPresent
		FROM attendanceData a
		JOIN meetingData m ON a.meetingID = m.meetingID
		JOIN courseFaculty cf ON m.classroomID = cf.classroomID
		JOIN courseData c ON cf.courseID = c.courseID
		WHERE a.studentID = $1
		GROUP BY c.courseID, c.courseName, cf.classroomID
	`

	rows, err := dbConn.Query(attendanceQuery, studentID)
	if err != nil {
		log.Printf("Error fetching attendance data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch attendance data"})
	}
	defer rows.Close()

	attendanceMap := make(map[int]models.CourseAttendance)

	for rows.Next() {
		var record models.CourseAttendance
		if err := rows.Scan(&record.CourseID, &record.CourseName, &record.ClassroomID, &record.TotalClasses, &record.NumPresent); err != nil {
			log.Printf("Error scanning attendance data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing attendance data"})
		}
		attendanceMap[record.CourseID] = record
	}

	response := models.StudentAttendanceResponse{
		StudentID:          studentID,
		AttendanceByCourse: attendanceMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}