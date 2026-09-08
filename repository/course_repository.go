package repository

import (
	"context"
	"fmt"
	"strings"

	"sql/database"
	"sql/models"
)

// CourseFilter holds optional filters for listing courses.
type CourseFilter struct {
	DepartmentID int
	Level        int
}

// GetCourses returns courses matching the given filters.
// Zero values mean "don't filter on this field".
func GetCourses(ctx context.Context, f CourseFilter) ([]models.Course, error) {
	conditions := []string{}
	args := []interface{}{}
	argPos := 1

	if f.DepartmentID != 0 {
		conditions = append(conditions, fmt.Sprintf("department_id = $%d", argPos))
		args = append(args, f.DepartmentID)
		argPos++
	}
	if f.Level != 0 {
		conditions = append(conditions, fmt.Sprintf("level = $%d", argPos))
		args = append(args, f.Level)
		argPos++
	}

	query := `SELECT course_id, department_id, course_code, course_title, credit_units, level, created_at FROM courses`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY course_id"

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		if err := rows.Scan(&c.CourseID, &c.DepartmentID, &c.CourseCode, &c.CourseTitle, &c.CreditUnits, &c.Level, &c.CreatedAt); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	return courses, rows.Err()
}