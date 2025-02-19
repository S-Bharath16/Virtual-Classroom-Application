package models

import "time"

type CourseFaculty struct {
	ClassroomID uint      `gorm:"primaryKey;autoIncrement" json:"classroomID"`
	CourseID    uint      `gorm:"not null" json:"courseID"`
	FacultyID   uint      `gorm:"not null" json:"facultyID"`
	SectionID   uint      `gorm:"not null" json:"sectionID"`
	SemesterID  uint      `gorm:"not null" json:"semesterID"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   *uint     `json:"createdBy"` // adminID
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedBy   *uint     `json:"updatedBy"` // adminID
}