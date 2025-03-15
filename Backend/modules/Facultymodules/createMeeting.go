package Facultymodules

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"Backend/database"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"golang.org/x/oauth2"
)

func CreateMeeting(c *fiber.Ctx) error {
    var request struct {
        ClassroomID       int       `json:"classroomID"`
		FacultyEmail	  string 	`json:"emailID"`
        StartTime         time.Time `json:"startTime"`
        EndTime           time.Time `json:"endTime"`
        MeetingDescription string   `json:"meetingDescription"`
        AccessToken       string    `json:"accessToken"`
    }

    if err := c.BodyParser(&request); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    dbConn, err := database.GetDB().DB()
    if err != nil {
        log.Printf("Error getting DB connection: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to get database connection",
        })
    }

    query := `
        SELECT cf.courseID, f.facultyName, s.emailID 
        FROM studentData s
        JOIN courseFaculty cf ON s.deptID = cf.deptID 
            AND s.semesterID = cf.semesterID 
            AND s.sectionID = cf.sectionID 
        JOIN facultyData f ON cf.facultyID = f.facultyID
        WHERE cf.classroomID = $1

        UNION 

        SELECT cf.courseID, f.facultyName, f.emailID 
        FROM facultyData f
        JOIN courseFaculty cf ON f.facultyID = cf.facultyID
        WHERE cf.classroomID = $1
    `

    rows, err := dbConn.Query(query, request.ClassroomID)
    if err != nil {
        log.Printf("Error fetching email IDs: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch email IDs",
        })
    }
    defer rows.Close()

    var courseID int
    var facultyName string
    attendees := []string{}

    for rows.Next() {
        var email string
        err := rows.Scan(&courseID, &facultyName, &email)
        if err != nil {
            log.Printf("Error scanning data: %v", err)
            continue
        }
        attendees = append(attendees, email)
    }

    facultyNameFormatted := strings.ReplaceAll(facultyName, " ", "")
    conferenceID := fmt.Sprintf("%d-%s-%d-%d", courseID, facultyNameFormatted, request.ClassroomID, time.Now().Unix())

    meetingLink, err := createGoogleMeetLink(conferenceID, attendees, request.StartTime, request.EndTime, request.MeetingDescription, request.AccessToken)
    if err != nil {
        log.Printf("Error generating Google Meet link: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to generate meeting link",
        })
    }

    insertQuery := `
        INSERT INTO meetingData (startTime, endTime, classroomID, meetingLink, createdBy, meetingDescription)
        VALUES ($1, $2, $3, $4, (SELECT facultyID FROM courseFaculty WHERE classroomID = $3), $5)
        RETURNING meetingID
    `
    var meetingID int
    err = dbConn.QueryRow(insertQuery, request.StartTime, request.EndTime, request.ClassroomID, meetingLink, request.MeetingDescription).Scan(&meetingID)
    if err != nil {
        log.Printf("Error inserting meeting data: %v", err)
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to insert meeting data",
        })
    }

    return c.Status(http.StatusOK).JSON(fiber.Map{
        "message":     "Meeting created successfully",
        "meetingID":   meetingID,
        "meetingLink": meetingLink,
    })
}

func createGoogleMeetLink(conferenceID string, attendees []string, startTime, endTime time.Time, description, accessToken string) (string, error) {
    ctx := context.Background()

    tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
    httpClient := oauth2.NewClient(ctx, tokenSource)

    service, err := calendar.NewService(ctx, option.WithHTTPClient(httpClient))
    if err != nil {
        return "", err
    }

    var eventAttendees []*calendar.EventAttendee
    for _, email := range attendees {
        eventAttendees = append(eventAttendees, &calendar.EventAttendee{Email: email})
    }

    event := &calendar.Event{
        Summary:     "Meeting",
        Description: description,
        Start:       &calendar.EventDateTime{DateTime: startTime.Format(time.RFC3339)},
        End:         &calendar.EventDateTime{DateTime: endTime.Format(time.RFC3339)},
        Attendees:   eventAttendees,
        ConferenceData: &calendar.ConferenceData{
            CreateRequest: &calendar.CreateConferenceRequest{
                RequestId: conferenceID,
                ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
                    Type: "hangoutsMeet",
                },
            },
        },
    }

    event, err = service.Events.Insert("primary", event).
        ConferenceDataVersion(1).
        SendUpdates("all").
        Do()
    if err != nil {
        return "", err
    }

    if event.ConferenceData == nil || len(event.ConferenceData.EntryPoints) == 0 {
        return "", fmt.Errorf("no conference data or entry points returned")
    }

    return event.ConferenceData.EntryPoints[0].Uri, nil
}