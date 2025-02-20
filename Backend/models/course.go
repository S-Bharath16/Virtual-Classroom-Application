package models

import "time"

type Course struct {
	CourseID     uint      `gorm:"primaryKey;autoIncrement" json:"courseID"`
	CourseCode   string    `gorm:"unique;not null" json:"courseCode"`
	CourseName   string    `gorm:"not null" json:"courseName"`
	CourseDeptID uint      `gorm:"not null" json:"courseDeptID"` 
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CourseType   string    `gorm:"not null" json:"courseType"` 
	UpdatedBy    *uint     `json:"updatedBy"`                  
}