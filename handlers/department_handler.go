package handlers

import (
	"strconv"

	"sql/repository"

	"github.com/gofiber/fiber/v3"
)

// GetDepartmentsHandler handles GET /departments
func GetDepartmentsHandler(c fiber.Ctx) error {
	facultyID := 0
	if id, err := strconv.Atoi(c.Query("faculty_id")); err == nil {
		facultyID = id
	}

	departments, err := repository.GetDepartments(c.Context(), facultyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": departments})
}
