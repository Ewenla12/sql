package repository

import (
	"context"

	"sql/database"
	"sql/models"
)

// GetFaculties returns every faculty in the system.
func GetFaculties(ctx context.Context) ([]models.Faculty, error) {
	rows, err := database.DB.Query(ctx, `SELECT faculty_id, faculty_name, created_at FROM faculties ORDER BY faculty_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faculties []models.Faculty
	for rows.Next() {
		var f models.Faculty
		if err := rows.Scan(&f.FacultyID, &f.FacultyName, &f.CreatedAt); err != nil {
			return nil, err
		}
		faculties = append(faculties, f)
	}

	return faculties, rows.Err()
}
