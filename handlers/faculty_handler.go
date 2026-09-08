package handlers

import (
	"sql/repository"

	"github.com/gofiber/fiber/v3"
)

// GetFacultiesHandler handles GET /faculties
func GetFacultiesHandler(c fiber.Ctx) error {
	faculties, err := repository.GetFaculties(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": faculties})
}
