package repository

import (
    "context"
    "errors"
    "fmt"
    "strings"

    "sql/database"
    "sql/models"

    "golang.org/x/crypto/bcrypt"
)

// CreateStudent hashes the given plain-text password, then inserts a new
// student row into Postgres. Returns an error if anything fails (e.g. the
// students_id or email already exists, violating a UNIQUE constraint).
func CreateStudent(ctx context.Context, s *models.Student, plainPassword string) error {
    if plainPassword == "" {
        return errors.New("password is required")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := `
        INSERT INTO students (
            students_id, first_name, last_name, date_of_birth,
            department_id, level, admission_year, email, password_hash
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
    `

    _, err = database.DB.Exec(ctx, query,
        s.StudentsID,
        s.FirstName,
        s.LastName,
        s.DateOfBirth,
        s.DepartmentID,
        s.Level,
        s.AdmissionYear,
        s.Email,
        string(hash),
    )

    return err
}

// StudentFilter holds all optional filters for listing students.
// Zero values (empty string / 0) mean "don't filter on this field".
type StudentFilter struct {
    DepartmentID int
    FacultyID    int
    Level        string
    FirstName    string
    LastName     string
    Status       string
    Page         int
    Limit        int
}

// GetStudents returns students matching the given filters, plus the total
// count of matching rows (before pagination) for building "page X of Y" UIs.
func GetStudents(ctx context.Context, f StudentFilter) ([]models.Student, int, error) {
    conditions := []string{}
    args := []interface{}{}
    argPos := 1

    if f.DepartmentID != 0 {
        conditions = append(conditions, fmt.Sprintf("s.department_id = $%d", argPos))
        args = append(args, f.DepartmentID)
        argPos++
    }
    if f.FacultyID != 0 {
        conditions = append(conditions, fmt.Sprintf("d.faculty_id = $%d", argPos))
        args = append(args, f.FacultyID)
        argPos++
    }
    if f.Level != "" {
        conditions = append(conditions, fmt.Sprintf("s.level = $%d", argPos))
        args = append(args, f.Level)
        argPos++
    }
    if f.FirstName != "" {
        conditions = append(conditions, fmt.Sprintf("s.first_name ILIKE $%d", argPos))
        args = append(args, "%"+f.FirstName+"%")
        argPos++
    }
    if f.LastName != "" {
        conditions = append(conditions, fmt.Sprintf("s.last_name ILIKE $%d", argPos))
        args = append(args, "%"+f.LastName+"%")
        argPos++
    }
    if f.Status != "" {
        conditions = append(conditions, fmt.Sprintf("s.status = $%d", argPos))
        args = append(args, f.Status)
        argPos++
    }

    whereClause := ""
    if len(conditions) > 0 {
        whereClause = "WHERE " + strings.Join(conditions, " AND ")
    }

    countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM students s
        JOIN departments d ON s.department_id = d.department_id
        %s
    `, whereClause)

    var total int
    if err := database.DB.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, err
    }

    if f.Page <= 0 {
        f.Page = 1
    }
    if f.Limit <= 0 {
        f.Limit = 20
    }
    offset := (f.Page - 1) * f.Limit

    dataQuery := fmt.Sprintf(`
        SELECT s.id, s.students_id, s.first_name, s.last_name, s.date_of_birth,
               s.department_id, s.level, s.status, s.admission_year, s.gender,
               s.email, s.password_hash, s.cgpa, s.created_at, s.updated_at
        FROM students s
        JOIN departments d ON s.department_id = d.department_id
        %s
        ORDER BY s.id
        LIMIT $%d OFFSET $%d
    `, whereClause, argPos, argPos+1)

    args = append(args, f.Limit, offset)

    rows, err := database.DB.Query(ctx, dataQuery, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var students []models.Student
    for rows.Next() {
        var s models.Student
        err := rows.Scan(
            &s.ID, &s.StudentsID, &s.FirstName, &s.LastName, &s.DateOfBirth,
            &s.DepartmentID, &s.Level, &s.Status, &s.AdmissionYear, &s.Gender,
            &s.Email, &s.PasswordHash, &s.CGPA, &s.CreatedAt, &s.UpdatedAt,
        )
        if err != nil {
            return nil, 0, err
        }
        students = append(students, s)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, err
    }

    return students, total, nil
}

// StudentUpdate holds optional fields for updating a student.
// Pointer fields let us distinguish "not provided" (nil) from
// "explicitly set to zero value" (e.g. an empty string).
type StudentUpdate struct {
    FirstName    *string
    LastName     *string
    DepartmentID *int
    Level        *string
    Email        *string
    Status       *string
    CGPA         *float64
}

// UpdateStudent updates only the fields that are non-nil in u.
// Returns an error if the student does not exist or the update fails.
func UpdateStudent(ctx context.Context, studentsID string, u StudentUpdate) error {
    setClauses := []string{}
    args := []interface{}{}
    argPos := 1

    if u.FirstName != nil {
        setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argPos))
        args = append(args, *u.FirstName)
        argPos++
    }
    if u.LastName != nil {
        setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argPos))
        args = append(args, *u.LastName)
        argPos++
    }
    if u.DepartmentID != nil {
        setClauses = append(setClauses, fmt.Sprintf("department_id = $%d", argPos))
        args = append(args, *u.DepartmentID)
        argPos++
    }
    if u.Level != nil {
        setClauses = append(setClauses, fmt.Sprintf("level = $%d", argPos))
        args = append(args, *u.Level)
        argPos++
    }
    if u.Email != nil {
        setClauses = append(setClauses, fmt.Sprintf("email = $%d", argPos))
        args = append(args, *u.Email)
        argPos++
    }
    if u.Status != nil {
        setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
        args = append(args, *u.Status)
        argPos++
    }
    if u.CGPA != nil {
        setClauses = append(setClauses, fmt.Sprintf("cgpa = $%d", argPos))
        args = append(args, *u.CGPA)
        argPos++
    }

    if len(setClauses) == 0 {
        return errors.New("no fields provided to update")
    }

    setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

    query := fmt.Sprintf(`
        UPDATE students
        SET %s
        WHERE students_id = $%d
    `, strings.Join(setClauses, ", "), argPos)

    args = append(args, studentsID)

    tag, err := database.DB.Exec(ctx, query, args...)
    if err != nil {
        return err
    }
    if tag.RowsAffected() == 0 {
        return errors.New("student not found")
    }

    return nil
}

// DeleteStudent removes a student row by their students_id (matric number).
func DeleteStudent(ctx context.Context, studentsID string) error {
    tag, err := database.DB.Exec(ctx, `DELETE FROM students WHERE students_id = $1`, studentsID)
    if err != nil {
        return err
    }
    if tag.RowsAffected() == 0 {
        return errors.New("student not found")
    }
    return nil
}
