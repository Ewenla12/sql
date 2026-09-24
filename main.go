package main

import (
	"log"
	"net"

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

	ports := []string{cfg.Port, "8080"}
	var listener net.Listener
	var err error
	var chosenPort string

	for _, port := range ports {
		if port == "" {
			continue
		}
		listener, err = net.Listen("tcp4", ":"+port)
		if err == nil {
			chosenPort = port
			break
		}
		log.Printf("Port %s unavailable: %v", port, err)
	}

	if listener == nil {
		log.Fatalf("failed to start server on any fallback port: %v", err)
	}

	log.Printf("Server listening on http://localhost:%s", chosenPort)
	log.Fatal(app.Listener(listener))
}
