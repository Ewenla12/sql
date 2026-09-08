package models

import "time"

type Department struct {
	DepartmentID   int       `json:"department_id"`
	FacultyID      int       `json:"faculty_id"`
	DepartmentName string    `json:"department_name"`
	DepartmentCode string    `json:"department_code"`
	CreatedAt      time.Time `json:"created_at"`
}
