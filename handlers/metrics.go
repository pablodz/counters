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
	userId := c.Get("user_id")

	if itemType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item_type required"})
	}
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item_id required"})
	}
	if eventType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "event_type required"})
	}
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id required"})
	}

	log := models.AuditLogPayload{
		UserID:    userId,
		UserType:  "registered",
		ItemID:    itemID,
		ItemType:  itemType,
		EventType: eventType,
		CreatedAt: time.Now().Unix(),
	}

	if err := store.LogInteraction(log); err != nil {
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
	}
	return c.JSON(result)
}
