package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/pablodz/counters/handlers"
	"github.com/pablodz/counters/singleton"
)

func main() {
	singleton.ValidateRequiredEnv()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${url}\n",
	}))

	app.Get("/healthz", handlers.Health)

	v1 := app.Group("/api/v1")
	v1.Post("/metrics/batch", handlers.BatchIncrement)
	v1.Post("/metrics/:content_type/:content_id/view", handlers.View())
	v1.Post("/metrics/:content_type/:content_id/like", handlers.Like())
	v1.Post("/metrics/:content_type/:content_id/share", handlers.Share())
	v1.Post("/metrics/:content_type/:content_id/increment", handlers.Increment)
	v1.Post("/metrics/:content_type/:content_id/reset", handlers.Reset)
	v1.Get("/metrics/:content_type/:content_id/:field", handlers.GetField)
	v1.Get("/metrics/:content_type/:content_id", handlers.GetMetrics)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
