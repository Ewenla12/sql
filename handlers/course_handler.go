package handlers

import (
	"strconv"

	"sql/repository"

	"github.com/gofiber/fiber/v3"
)

// GetCoursesHandler handles GET /courses
func GetCoursesHandler(c fiber.Ctx) error {
	filter := repository.CourseFilter{}

	if deptID, err := strconv.Atoi(c.Query("department_id")); err == nil {
		filter.DepartmentID = deptID
	}
	if level, err := strconv.Atoi(c.Query("level")); err == nil {
		filter.Level = level
	}

	courses, err := repository.GetCourses(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": courses})
}
