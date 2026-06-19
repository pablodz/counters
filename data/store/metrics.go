package store

import (
	"encoding/json"
	"fmt"

	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/singleton"
)

const selectColumns = "content_id, content_type, views_count, likes_count, shares_count, updated_at"

func GetMetrics(contentType, contentID string) (*models.Metrics, error) {
	raw, err := singleton.D1Exec(
		"SELECT "+selectColumns+" FROM post_metrics WHERE content_id = ? AND content_type = ?",
		contentID, contentType,
	)
	if err != nil {
		return nil, err
	}
	var rows []models.Metrics
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &models.Metrics{ContentID: contentID, ContentType: contentType}, nil
	}
	return &rows[0], nil
}

func IncrementMetric(contentType, contentID, field string, amount int) (*models.Metrics, error) {
	views, likes, shares := 0, 0, 0
	switch field {
	case "views_count":
		views = amount
	case "likes_count":
		likes = amount
	case "shares_count":
		shares = amount
	default:
		return nil, fmt.Errorf("invalid field: %s", field)
	}

	sql := `INSERT INTO post_metrics (content_id, content_type, views_count, likes_count, shares_count, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(content_id, content_type) DO UPDATE SET
			views_count = views_count + excluded.views_count,
			likes_count = likes_count + excluded.likes_count,
			shares_count = shares_count + excluded.shares_count,
			updated_at = CURRENT_TIMESTAMP`
	if _, err := singleton.D1Exec(sql, contentID, contentType, views, likes, shares); err != nil {
		return nil, err
	}
	return GetMetrics(contentType, contentID)
}

func ResetMetric(contentType, contentID string) (*models.Metrics, error) {
	if _, err := singleton.D1Exec(
		"UPDATE post_metrics SET views_count = 0, likes_count = 0, shares_count = 0, updated_at = CURRENT_TIMESTAMP WHERE content_id = ? AND content_type = ?",
		contentID, contentType,
	); err != nil {
		return nil, err
	}
	return GetMetrics(contentType, contentID)
}
