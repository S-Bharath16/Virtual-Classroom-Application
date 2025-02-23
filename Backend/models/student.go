package models

type Student struct {
	StudentID   uint   `gorm:"primaryKey;autoIncrement" json:"studentID"`
	RollNumber  string `gorm:"unique;not null" json:"rollNumber"`
	EmailID     string `gorm:"unique;not null" json:"emailId"`
	StudentName string `gorm:"not null" json:"studentName"`
	StartYear   int    `gorm:"not null" json:"startYear"`
	EndYear     int    `gorm:"not null" json:"endYear"`
	DeptID      *uint  `gorm:"foreignKey:DeptID;references:DeptID;onDelete:SET NULL" json:"deptID"`
	SectionID   *uint  `gorm:"foreignKey:SectionID;references:SectionID;onDelete:SET NULL" json:"sectionID"`
	SemesterID  *uint  `gorm:"foreignKey:SemesterID;references:SemesterID;onDelete:SET NULL" json:"semesterID"`
}

type StudentProfile struct {
	RollNumber     string  `json:"rollNumber"`
	EmailID        string  `json:"emailID"`
	StudentName    string  `json:"studentName"`
	StartYear      int     `json:"startYear"`
	EndYear        int     `json:"endYear"`
	DeptName       *string `json:"deptName"`       // from deptData.deptName
	SectionName    *string `json:"sectionName"`    // from sectionData.sectionName
	SemesterNumber *int    `json:"semesterNumber"` // from semesterData.semesterNumber
}