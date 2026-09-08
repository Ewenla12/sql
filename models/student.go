package models

import "time"

type Student struct {
	ID            int           `json:"id"`
	StudentsID    string        `json:"students_id"` // matric number, public-facing
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	DateOfBirth   *time.Time    `json:"date_of_birth,omitempty"`
	DepartmentID  int           `json:"department_id"`
	Level         StudentLevel  `json:"level"`
	Status        StudentStatus `json:"status"`
	AdmissionYear int           `json:"admission_year"`
	Gender        *string       `json:"gender,omitempty"`
	Email         *string       `json:"email,omitempty"`
	PasswordHash  string        `json:"-"`
	CGPA          float64       `json:"cgpa"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
