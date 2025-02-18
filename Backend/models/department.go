package models

type Department struct{
	DepartmentID   uint   `gorm:"primaryKey;autoIncrement" json:"deptID"`
	DepartmentName  string `gorm:"unique;not null" json:"deptName"`
}