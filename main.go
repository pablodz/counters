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

	v1 := app.Group("/api/v1")
	v1.Get("/:item_type/:item_id/:event_type/:user_id", handlers.IncrementEvent)
	v1.Get("/:item_type/:item_id", handlers.GetMetrics)
	v1.Get("/histogram/:item_type/:item_id", handlers.GetHistogram)
	v1.Get("/activity/:item_type/:item_id", handlers.GetRecentActivity)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
