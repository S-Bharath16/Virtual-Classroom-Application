package models

import "time"

type Assignment struct {
	AssignmentID        uint      `gorm:"primaryKey;autoIncrement" json:"assignmentID"`
	ClassroomID         uint      `gorm:"not null" json:"classroomID"`
	AssignmentName      string    `gorm:"not null" json:"assignmentName"`
	AssignmentDescription string `json:"assignmentDescription"`
	AssignmentText      string    `json:"assignmentText"`
	ImagePath           *string   `json:"imagePath"`
	DocumentPath        *string   `json:"documentPath"`
	StartTime           time.Time `json:"startTime"`
	EndTime             time.Time `json:"endTime"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	CreatedBy           uint      `json:"createdBy"`
}
