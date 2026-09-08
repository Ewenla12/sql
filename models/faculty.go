package models

import "time"

type Faculty struct {
	FacultyID   int       `json:"faculty_id"`
	FacultyName string    `json:"faculty_name"`
	CreatedAt   time.Time `json:"created_at"`
}
