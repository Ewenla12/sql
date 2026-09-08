package models

import "time"

type Course struct {
	CourseID     int       `json:"course_id"`
	DepartmentID int       `json:"department_id"`
	CourseCode   string    `json:"course_code"`
	CourseTitle  string    `json:"course_title"`
	CreditUnits  int       `json:"credit_units"`
	Level        int       `json:"level"`
	CreatedAt    time.Time `json:"created_at"`
}
