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

	app.Get(handlers.IncrementEventPattern, handlers.IncrementEvent)
	app.Get(handlers.GetMetricsPattern, handlers.GetMetrics)
	app.Get(handlers.GetMetricsListPattern, handlers.GetMetricsList)
	app.Get(handlers.GetHistogramPattern, handlers.GetHistogram)
	app.Get(handlers.GetRecentActivityPattern, handlers.GetRecentActivity)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
