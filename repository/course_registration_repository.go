package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"sql/database"
	"sql/models"
)

// ApplyForCourses inserts one pending registration row per course_id.
// Returns the list of new registration_ids created.
func ApplyForCourses(ctx context.Context, studentID string, academicSession string, courseIDs []int) ([]int, error) {
	if len(courseIDs) == 0 {
		return nil, errors.New("at least one course_id is required")
	}

	var registrationIDs []int

	for _, courseID := range courseIDs {
		var id int
		err := database.DB.QueryRow(ctx, `
			INSERT INTO course_registration (student_id, course_id, academic_session, status)
			VALUES ($1, $2, $3, 'pending')
			RETURNING registration_id
		`, studentID, courseID, academicSession).Scan(&id)

		if err != nil {
			return nil, err
		}
		registrationIDs = append(registrationIDs, id)
	}

	return registrationIDs, nil
}

// CourseRegistrationFilter holds optional filters for admin listing.
type CourseRegistrationFilter struct {
	Status          string
	DepartmentID    int
	AcademicSession string
}

// CourseRegistrationView is a joined, read-friendly shape for the admin
// listing endpoint (student + course details alongside the registration).
type CourseRegistrationView struct {
	RegistrationID  int       `json:"registration_id"`
	AcademicSession string    `json:"academic_session"`
	Status          string    `json:"status"`
	AppliedAt       time.Time `json:"applied_at"`
	StudentID       string    `json:"student_id"`
	StudentFirst    string    `json:"student_first_name"`
	StudentLast     string    `json:"student_last_name"`
	StudentLevel    string    `json:"student_level"`
	CourseID        int       `json:"course_id"`
	CourseCode      string    `json:"course_code"`
	CourseTitle     string    `json:"course_title"`
	CourseUnits     int       `json:"course_credit_units"`
}

// GetCourseRegistrations returns registrations matching the filter, joined
// with student and course info for a rich admin view.
func GetCourseRegistrations(ctx context.Context, f CourseRegistrationFilter) ([]CourseRegistrationView, error) {
	conditions := []string{}
	args := []interface{}{}
	argPos := 1

	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("cr.status = $%d", argPos))
		args = append(args, f.Status)
		argPos++
	}
	if f.DepartmentID != 0 {
		conditions = append(conditions, fmt.Sprintf("s.department_id = $%d", argPos))
		args = append(args, f.DepartmentID)
		argPos++
	}
	if f.AcademicSession != "" {
		conditions = append(conditions, fmt.Sprintf("cr.academic_session = $%d", argPos))
		args = append(args, f.AcademicSession)
		argPos++
	}

	query := `
		SELECT cr.registration_id, cr.academic_session, cr.status, cr.applied_at,
		       s.students_id, s.first_name, s.last_name, s.level,
		       c.course_id, c.course_code, c.course_title, c.credit_units
		FROM course_registration cr
		JOIN students s ON cr.student_id = s.students_id
		JOIN courses c ON cr.course_id = c.course_id
	`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY cr.applied_at DESC"

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CourseRegistrationView
	for rows.Next() {
		var v CourseRegistrationView
		err := rows.Scan(
			&v.RegistrationID, &v.AcademicSession, &v.Status, &v.AppliedAt,
			&v.StudentID, &v.StudentFirst, &v.StudentLast, &v.StudentLevel,
			&v.CourseID, &v.CourseCode, &v.CourseTitle, &v.CourseUnits,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}

	return results, rows.Err()
}

// ApproveRegistration marks a registration as approved and creates the
// corresponding student_courses row, all in a single transaction — if
// either step fails, neither change is saved.
func ApproveRegistration(ctx context.Context, registrationID int, reviewedByAdminID int) error {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op if already committed

	var studentID string
	var courseID int
	var academicSession string

	err = tx.QueryRow(ctx, `
		SELECT student_id, course_id, academic_session
		FROM course_registration
		WHERE registration_id = $1 AND status = 'pending'
	`, registrationID).Scan(&studentID, &courseID, &academicSession)
	if err != nil {
		return errors.New("pending registration not found")
	}

	_, err = tx.Exec(ctx, `
		UPDATE course_registration
		SET status = 'approved', reviewed_at = CURRENT_TIMESTAMP, reviewed_by = $1
		WHERE registration_id = $2
	`, reviewedByAdminID, registrationID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO student_courses (student_id, course_id, academic_session)
		VALUES ($1, $2, $3)
	`, studentID, courseID, academicSession)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// RejectRegistration marks a registration as rejected.
func RejectRegistration(ctx context.Context, registrationID int, reviewedByAdminID int) error {
	tag, err := database.DB.Exec(ctx, `
		UPDATE course_registration
		SET status = 'rejected', reviewed_at = CURRENT_TIMESTAMP, reviewed_by = $1
		WHERE registration_id = $2 AND status = 'pending'
	`, reviewedByAdminID, registrationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("pending registration not found")
	}
	return nil
}

// GetStudentCourses returns a student's approved/assigned courses for a session.
func GetStudentCourses(ctx context.Context, studentID string, academicSession string) ([]models.StudentCourse, error) {
	query := `
		SELECT id, student_id, course_id, academic_session, score, grade, assigned_at
		FROM student_courses
		WHERE student_id = $1
	`
	args := []interface{}{studentID}

	if academicSession != "" {
		query += " AND academic_session = $2"
		args = append(args, academicSession)
	}
	query += " ORDER BY id"

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.StudentCourse
	for rows.Next() {
		var sc models.StudentCourse
		if err := rows.Scan(&sc.ID, &sc.StudentID, &sc.CourseID, &sc.AcademicSession, &sc.Score, &sc.Grade, &sc.AssignedAt); err != nil {
			return nil, err
		}
		courses = append(courses, sc)
	}

	return courses, rows.Err()
}
