package middleware

import (
	"strings"

	"sql/utils"

	"github.com/gofiber/fiber/v3"
)

// RequireAuth is middleware that checks for a valid JWT in the
// Authorization header. If valid, it stores the claims on the context
// (accessible via c.Locals("claims")) so downstream handlers know who's
// making the request. If invalid or missing, it stops the request early
// with 401 Unauthorized.
func RequireAuth(c fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
	}

	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	c.Locals("claims", claims)

	return c.Next()
}

// RequireRole is middleware that must run AFTER RequireAuth. It checks
// that the authenticated user's role is one of the allowed roles.
// Use this to restrict admin-only endpoints to specific admin roles.
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*utils.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not authenticated"})
		}

		for _, role := range allowedRoles {
			if claims.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})
	}
}
