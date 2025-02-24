package mailer

import (
	"log"
	"time"
	"errors"
	"Backend/database"
	"Backend/utilities/mailer"
)

func QuizCreationMail(Subject string, Body string, ClassroomID int) error {

	dbConn, err := database.GetDB().DB()
	if err != nil {
		log.Printf("[ERROR]: Error getting DB connection: %v", err)
		return errors.New("failed to get database connection")
	}

	var deptID, semesterID, sectionID, courseID int
	query := `
		SELECT deptID, semesterID, sectionID, courseID
		FROM courseFaculty 
		WHERE classroomID = $1
	`
	err = dbConn.QueryRow(query, ClassroomID).Scan(&deptID, &semesterID, &sectionID, &courseID)
	if err != nil {
		log.Printf("[ERROR]: Failed to fetch classroom details: %v", err)
		return errors.New("failed to fetch classroom details")
	}

	var courseName string
	courseQuery := `
		SELECT courseName FROM courseData 
		WHERE courseID = $1
	`
	err = dbConn.QueryRow(courseQuery, courseID).Scan(&courseName)
	if err != nil {
		log.Printf("[ERROR]: Failed to fetch course name: %v", err)
		return errors.New("failed to fetch course name")
	}

	var startTime, endTime time.Time
	quizQuery := `
		SELECT startTime, endTime FROM quizData
		WHERE classroomID = $1
	`
	err = dbConn.QueryRow(quizQuery, ClassroomID).Scan(&startTime, &endTime)
	if err != nil {
		log.Printf("[ERROR]: Failed to fetch quiz timing: %v", err)
		return errors.New("failed to fetch quiz timing")
	}

	startTimeFormatted := startTime.Format("2006-01-02 15:04:05")
	endTimeFormatted := endTime.Format("2006-01-02 15:04:05")

	emailQuery := `
		SELECT emailID FROM studentData 
		WHERE deptID = $1 AND semesterID = $2 AND sectionID = $3
	`
	rows, err := dbConn.Query(emailQuery, deptID, semesterID, sectionID)
	if err != nil {
		log.Printf("[ERROR]: Failed to fetch student emails: %v", err)
		return errors.New("failed to fetch student emails")
	}
	defer rows.Close()

	var emailIDs []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Printf("[ERROR]: Failed to scan email: %v", err)
			continue
		}
		emailIDs = append(emailIDs, email)
	}

	if len(emailIDs) == 0 {
		log.Println("[INFO]: No students found for the given classroom.")
		return errors.New("no students found for the given classroom")
	}

	err = mailer.SendMail(emailIDs, Subject, Body, courseName, startTimeFormatted, endTimeFormatted, "quizCreation")
	if err != nil {
		log.Printf("[ERROR]: Failed to send emails: %v", err)
		return errors.New("failed to send emails")
	}

	log.Println("[LOG]: Quiz creation emails sent successfully")
	return nil
}