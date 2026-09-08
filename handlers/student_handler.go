package handlers

import (
	"strconv"
	"time"

	"sql/models"
	"sql/repository"

	"github.com/gofiber/fiber/v3"
)

// CreateStudentHandler handles POST /students
func CreateStudentHandler(c fiber.Ctx) error {
	type request struct {
		StudentID     string `json:"student_id"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		DateOfBirth   string `json:"date_of_birth"`
		DepartmentID  int    `json:"department_id"`
		Level         string `json:"level"`
		AdmissionYear int    `json:"admission_year"`
		Email         string `json:"email"`
		Password      string `json:"password"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.StudentID == "" || req.FirstName == "" || req.LastName == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "student_id, first_name, last_name, and password are required"})
	}

	student := &models.Student{
		StudentsID:    req.StudentID,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		DepartmentID:  req.DepartmentID,
		Level:         models.StudentLevel(req.Level),
		AdmissionYear: req.AdmissionYear,
		Email:         &req.Email,
	}

	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "date_of_birth must be in YYYY-MM-DD format"})
		}
		student.DateOfBirth = &dob
	}

	if err := repository.CreateStudent(c.Context(), student, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"student_id": req.StudentID,
		"status":     "active",
		"message":    "Student admitted successfully",
	})
}

// GetStudentsHandler handles GET /students
func GetStudentsHandler(c fiber.Ctx) error {
	filter := repository.StudentFilter{
		Level:     c.Query("level"),
		FirstName: c.Query("first_name"),
		LastName:  c.Query("last_name"),
		Status:    c.Query("status"),
	}

	if deptID, err := strconv.Atoi(c.Query("department_id")); err == nil {
		filter.DepartmentID = deptID
	}
	if facID, err := strconv.Atoi(c.Query("faculty_id")); err == nil {
		filter.FacultyID = facID
	}
	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		filter.Page = page
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		filter.Limit = limit
	}

	students, total, err := repository.GetStudents(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  students,
		"total": total,
		"page":  filter.Page,
	})
}

// GetStudentByIDHandler handles GET /students/:student_id
func GetStudentByIDHandler(c fiber.Ctx) error {
	studentID := c.Params("student_id")

	student, err := repository.GetStudentByID(c.Context(), studentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "student not found"})
	}

	return c.JSON(student)
}

// UpdateStudentHandler handles PUT /students/:student_id
func UpdateStudentHandler(c fiber.Ctx) error {
	studentID := c.Params("student_id")

	type request struct {
		FirstName    *string  `json:"first_name"`
		LastName     *string  `json:"last_name"`
		DepartmentID *int     `json:"department_id"`
		Level        *string  `json:"level"`
		Email        *string  `json:"email"`
		Status       *string  `json:"status"`
		CGPA         *float64 `json:"cgpa"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	update := repository.StudentUpdate{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		DepartmentID: req.DepartmentID,
		Level:        req.Level,
		Email:        req.Email,
		Status:       req.Status,
		CGPA:         req.CGPA,
	}

	if err := repository.UpdateStudent(c.Context(), studentID, update); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"student_id": studentID,
		"message":    "Student updated successfully",
	})
}

// DeleteStudentHandler handles DELETE /students/:student_id
func DeleteStudentHandler(c fiber.Ctx) error {
	studentID := c.Params("student_id")

	if err := repository.DeleteStudent(c.Context(), studentID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Student record deleted"})
}
