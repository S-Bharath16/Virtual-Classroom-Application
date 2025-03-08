package models

type Attendance struct {
	MeetingID int    `json:"meetingID"`
	StudentID int    `json:"studentID"`
	IsPresent string `json:"isPresent"`
}