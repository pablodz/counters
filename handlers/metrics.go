package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/data/store"
)

func Health(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

type incrementRequest struct {
	Field  string `json:"field"`
	Amount int    `json:"amount"`
}

func Increment(c fiber.Ctx) error {
	contentType := c.Params("content_type")
	contentID := c.Params("content_id")

	var body incrementRequest
	_ = c.Bind().Body(&body)
	if body.Field == "" {
		body.Field = "views_count"
	}
	if body.Amount == 0 {
		body.Amount = 1
	}
	if !models.ValidField(body.Field) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid field, must be one of: views_count, likes_count, shares_count"})
	}

	m, err := store.IncrementMetric(contentType, contentID, body.Field, body.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(m)
}

func shorthandIncrement(field string) fiber.Handler {
	return func(c fiber.Ctx) error {
		m, err := store.IncrementMetric(c.Params("content_type"), c.Params("content_id"), field, 1)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(m)
	}
}

func View() fiber.Handler  { return shorthandIncrement("views_count") }
func Like() fiber.Handler  { return shorthandIncrement("likes_count") }
func Share() fiber.Handler { return shorthandIncrement("shares_count") }

func GetMetrics(c fiber.Ctx) error {
	m, err := store.GetMetrics(c.Params("content_type"), c.Params("content_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(m)
}

func GetField(c fiber.Ctx) error {
	field := c.Params("field")
	if !models.ValidField(field) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid field"})
	}
	m, err := store.GetMetrics(c.Params("content_type"), c.Params("content_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	var val int
	switch field {
	case "views_count":
		val = m.ViewsCount
	case "likes_count":
		val = m.LikesCount
	case "shares_count":
		val = m.SharesCount
	}
	return c.JSON(fiber.Map{
		"content_id":   m.ContentID,
		"content_type": m.ContentType,
		field:          val,
	})
}

type batchEvent struct {
	ContentType string `json:"content_type"`
	ContentID   string `json:"content_id"`
	Field       string `json:"field"`
	Amount      int    `json:"amount"`
}

type batchRequest struct {
	Events []batchEvent `json:"events"`
}

func BatchIncrement(c fiber.Ctx) error {
	var body batchRequest
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	if len(body.Events) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "events array is empty"})
	}

	succeeded := 0
	var failed []map[string]string
	for _, e := range body.Events {
		if e.Field == "" {
			e.Field = "views_count"
		}
		if e.Amount == 0 {
			e.Amount = 1
		}
		if e.ContentType == "" || e.ContentID == "" {
			failed = append(failed, map[string]string{"content_id": e.ContentID, "error": "content_type and content_id required"})
			continue
		}
		if !models.ValidField(e.Field) {
			failed = append(failed, map[string]string{"content_id": e.ContentID, "error": "invalid field"})
			continue
		}
		if _, err := store.IncrementMetric(e.ContentType, e.ContentID, e.Field, e.Amount); err != nil {
			failed = append(failed, map[string]string{"content_id": e.ContentID, "error": err.Error()})
			continue
		}
		succeeded++
	}

	return c.JSON(fiber.Map{
		"succeeded": succeeded,
		"failed":    failed,
	})
}

func Reset(c fiber.Ctx) error {
	m, err := store.ResetMetric(c.Params("content_type"), c.Params("content_id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(m)
}
