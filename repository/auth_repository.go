package repository

import (
	"context"

	"sql/database"
	"sql/models"
)

// GetStudentByID fetches a single student by their students_id (matric number).
func GetStudentByID(ctx context.Context, studentsID string) (*models.Student, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT id, students_id, first_name, last_name, date_of_birth,
		       department_id, level, status, admission_year, gender,
		       email, password_hash, cgpa, created_at, updated_at
		FROM students
		WHERE students_id = $1
	`, studentsID)

	var s models.Student
	err := row.Scan(
		&s.ID, &s.StudentsID, &s.FirstName, &s.LastName, &s.DateOfBirth,
		&s.DepartmentID, &s.Level, &s.Status, &s.AdmissionYear, &s.Gender,
		&s.Email, &s.PasswordHash, &s.CGPA, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetAdminByUsername fetches a single admin by username.
func GetAdminByUsername(ctx context.Context, username string) (*models.Admin, error) {
	row := database.DB.QueryRow(ctx, `
		SELECT admin_id, username, email, password_hash, role, department_id, created_at
		FROM admin
		WHERE username = $1
	`, username)

	var a models.Admin
	err := row.Scan(
		&a.AdminID, &a.Username, &a.Email, &a.PasswordHash, &a.Role, &a.DepartmentID, &a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
