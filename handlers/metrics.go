package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/data/store"
)

func Health(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func GetMetrics(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")

	m, err := store.GetMetrics(itemID, itemType)
	if err != nil || m == nil {
		return c.JSON(&models.Metrics{
			ItemID:      itemID,
			ItemType:    itemType,
			ViewsCount:  0,
			LikesCount:  0,
			SharesCount: 0,
			UpdatedAt:   time.Now().Unix(),
		})
	}
	return c.JSON(m)
}

type eventRequest struct {
	ItemID    string `json:"item_id"`
	ItemType  string `json:"item_type"`
	EventType string `json:"event_type"`
}

func IncrementEvent(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")
	eventType := c.Params("event_type")

	if itemType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item_type required"})
	}
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "item_id required"})
	}
	if eventType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "event_type required"})
	}

	unixHour, err := models.PrepararDatosInteraccion(models.TrackingPayload{
		ItemID:    itemID,
		ItemType:  itemType,
		EventType: eventType,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := store.LogInteraction(itemID, itemType, eventType, unixHour); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

func GetHistogram(c fiber.Ctx) error {
	itemType := c.Params("item_type")
	itemID := c.Params("item_id")
	eventType := c.Params("event_type")

	resolution := c.Query("resolution", "1h")
	if _, ok := models.ResolutionSeconds[resolution]; !ok {
		resolution = "1h"
	}

	var from, to int64
	if v := c.Query("from"); v != "" {
		fmt.Sscanf(v, "%d", &from)
	}
	if v := c.Query("to"); v != "" {
		fmt.Sscanf(v, "%d", &to)
	}

	result, err := store.GetHistogram(itemID, itemType, eventType, resolution, from, to)
	if err != nil {
		return c.JSON([]models.HistogramBucket{})
	}
	return c.JSON(result)
}
