package models

type Student struct {
	StudentID   uint   `gorm:"primaryKey;autoIncrement" json:"studentID"`
	RollNumber  string `gorm:"unique;not null" json:"rollNumber"`
	EmailID     string `gorm:"unique;not null" json:"emailId"`
	StudentName string `gorm:"not null" json:"studentName"`
	StartYear   int    `gorm:"not null" json:"startYear"`
	EndYear     int    `gorm:"not null" json:"endYear"`
	DeptID      *uint  `gorm:"foreignKey:DeptID;references:DeptID;onDelete:SET NULL" json:"deptID"`
	Section     string `gorm:"not null" json:"section"`
	Semester    int    `gorm:"not null" json:"semester"`
}
