package models

type StudentCourse struct {
	ClassroomID uint   `json:"classroomID"`
	CourseID    uint   `json:"courseID"`
	CourseCode  string `json:"courseCode"`
	CourseName  string `json:"courseName"`
	FacultyID   uint   `json:"facultyID"`
	FacultyName string `json:"facultyName"`
}
