package handlers

import (
	"log"
	"strings"
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
	if err != nil {
		log.Printf("error getting metrics: %v", err)
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

func GetMetricsList(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemIDsParam := c.Query("item_ids", "")

	if itemType == "" || itemIDsParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item_type path param and item_ids query param are required"})
	}

	itemIDs := strings.Split(itemIDsParam, ",")

	result, err := store.GetMetricsList(itemType, itemIDs)
	if err != nil {
		log.Printf("error getting metrics list: %v", err)
	}

	return c.JSON(result)
}

func GetRecentActivity(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")

	result, err := store.GetRecentActivity(itemID, itemType)
	if err != nil {
		log.Printf("error getting recent activity: %v", err)
	}
	return c.JSON(result)
}
