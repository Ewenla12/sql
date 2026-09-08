package handlers

import (
	"strconv"

	"sql/repository"
	"sql/utils"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

// StudentLoginHandler handles POST /auth/student/login
func StudentLoginHandler(c fiber.Ctx) error {
	type request struct {
		StudentID string `json:"student_id"`
		Password  string `json:"password"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	student, err := repository.GetStudentByID(c.Context(), req.StudentID)
	if err != nil {
		// Deliberately vague: don't reveal whether the ID exists or the password was wrong.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := utils.GenerateToken(student.StudentsID, "student")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"student": fiber.Map{
			"student_id":    student.StudentsID,
			"first_name":    student.FirstName,
			"last_name":     student.LastName,
			"department_id": student.DepartmentID,
			"level":         student.Level,
		},
	})
}

// AdminLoginHandler handles POST /auth/admin/login
func AdminLoginHandler(c fiber.Ctx) error {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	admin, err := repository.GetAdminByUsername(c.Context(), req.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := utils.GenerateToken(strconv.Itoa(admin.AdminID), string(admin.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"admin": fiber.Map{
			"admin_id": admin.AdminID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
}
