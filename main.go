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
	v1.Get(handlers.IncrementEventPattern, handlers.IncrementEvent)
	v1.Get(handlers.GetMetricsPattern, handlers.GetMetrics)
	v1.Get(handlers.GetMetricsListPattern, handlers.GetMetricsList)
	v1.Get(handlers.GetHistogramPattern, handlers.GetHistogram)
	v1.Get(handlers.GetRecentActivityPattern, handlers.GetRecentActivity)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
