package models

import "time"

type Lecturer struct {
	LecturerID   int       `json:"lecturer_id"`
	StaffID      string    `json:"staff_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	DepartmentID int       `json:"department_id"`
	Email        *string   `json:"email,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
