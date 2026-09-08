package models

import "time"

type CourseRegistration struct {
	RegistrationID  int               `json:"registration_id"`
	StudentID       string            `json:"student_id"`
	CourseID        int               `json:"course_id"`
	AcademicSession string            `json:"academic_session"`
	AppliedAt       time.Time         `json:"applied_at"`
	ReviewedAt      *time.Time        `json:"reviewed_at,omitempty"`
	ReviewedBy      *int              `json:"reviewed_by,omitempty"`
	Status          ApplicationStatus `json:"status"`
}
