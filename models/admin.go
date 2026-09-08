package models

import "time"

type Admin struct {
	AdminID      int       `json:"admin_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         AdminRole `json:"role"`
	DepartmentID *int      `json:"department_id,omitempty"` // nullable — HOD tied to a dept, super_admin isn't
	CreatedAt    time.Time `json:"created_at"`
}
