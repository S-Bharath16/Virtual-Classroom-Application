package models

import "time"

type Faculty struct {
	FacultyID   uint      `gorm:"primaryKey;autoIncrement" json:"facultyID"`
	EmailID     string    `gorm:"unique;not null" json:"emailID"`
	FacultyName string    `gorm:"not null" json:"facultyName"`
	DeptID      *uint     `gorm:"foreignKey:DeptID;references:DeptID;onDelete:SET NULL" json:"deptID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}