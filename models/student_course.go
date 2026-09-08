package models

import "time"

type StudentCourse struct {
	ID              int       `json:"id"`
	StudentID       string    `json:"student_id"` // FK -> students.students_id
	CourseID        int       `json:"course_id"`
	AcademicSession string    `json:"academic_session"`
	Score           *float64  `json:"score,omitempty"`
	Grade           *string   `json:"grade,omitempty"`
	AssignedAt      time.Time `json:"assigned_at"`
}
