package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/data/store"
)

func IncrementEvent(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")
	eventType := c.Params("event_type")
	userId := c.Params("user_id")

	if itemType == "" || itemID == "" || eventType == "" || userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "all path parameters are required"})
	}

	logData := models.AuditLogPayload{
		UserID:    userId,
		UserType:  "registered",
		ItemID:    itemID,
		ItemType:  itemType,
		EventType: eventType,
		CreatedAt: time.Now().Unix(),
	}

	if err := store.LogInteraction(logData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

func GetMetrics(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")

	m, err := store.GetMetrics(itemID, itemType)
	if err != nil || m == nil {
		return c.JSON(map[string]int{
			"view":  0,
			"like":  0,
			"share": 0,
		})
	}

	return c.JSON(m)
}

func GetHistogram(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")
	resolution := c.Query("resolution", "1h")

	if !models.IsValidResolution(resolution) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid resolution"})
	}

	result, err := store.GetHistogram(itemID, itemType, resolution)
	if err != nil {
		log.Printf("error getting histogram: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate histogram"})
	}
	return c.JSON(result)
}
