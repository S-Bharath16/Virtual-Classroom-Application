package models

import "time"

type Quiz struct {
	QuizID         uint      `gorm:"primaryKey;autoIncrement" json:"quizID"`
	ClassroomID    uint      `json:"classroomID"`    
	QuizName       string    `gorm:"not null" json:"quizName"`
	QuizDescription string   `json:"quizDescription"` 
	QuizData       string    `json:"quizData"`        
	IsOpenForAll   string    `gorm:"default:'0'" json:"isOpenForAll"` 
	StartTime      time.Time `gorm:"not null" json:"startTime"`
	EndTime        time.Time `gorm:"not null" json:"endTime"`
	QuizDuration   int       `gorm:"not null" json:"quizDuration"` 
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
	CreatedBy      uint      `json:"createdBy"`      
}
