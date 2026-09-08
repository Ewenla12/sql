package main

import (
	"log"

	"sql/config"
	"sql/database"
	"sql/routes"
	"sql/utils"

	"github.com/gofiber/fiber/v3"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg.DatabaseURL)
	utils.JWTSecret = cfg.JWTSecret

	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	routes.RegisterRoutes(app)

	log.Fatal(app.Listen(":" + cfg.Port))
}
