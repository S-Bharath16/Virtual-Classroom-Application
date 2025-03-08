package models

type CourseAttendance struct {
	CourseID     int    `json:"courseID"`
	CourseName   string `json:"courseName"`
	ClassroomID  int    `json:"classroomID"`
	TotalClasses int    `json:"totalClasses"`
	NumPresent   int    `json:"numPresent"`
}

type StudentAttendanceResponse struct {
	StudentID          int                          `json:"studentID"`
	AttendanceByCourse map[int]CourseAttendance `json:"attendanceByCourse"`
}