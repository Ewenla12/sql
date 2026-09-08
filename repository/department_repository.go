package repository

import (
	"context"

	"sql/database"
	"sql/models"
)

// GetDepartments returns departments, optionally filtered by faculty_id.
// Pass facultyID = 0 to get all departments.
func GetDepartments(ctx context.Context, facultyID int) ([]models.Department, error) {
	query := `SELECT department_id, faculty_id, department_name, department_code, created_at FROM departments`
	args := []interface{}{}

	if facultyID != 0 {
		query += ` WHERE faculty_id = $1`
		args = append(args, facultyID)
	}
	query += ` ORDER BY department_id`

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.DepartmentID, &d.FacultyID, &d.DepartmentName, &d.DepartmentCode, &d.CreatedAt); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}

	return departments, rows.Err()
}
