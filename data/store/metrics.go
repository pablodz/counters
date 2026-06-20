package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/singleton"
)

func LogInteraction(itemID, itemType, eventType string, unixHour int64) error {
	sql := `INSERT INTO item_snapshots (item_id, item_type, event_type, period_hour_unix, total_count)
		VALUES (?, ?, ?, ?, 1)
		ON CONFLICT(item_id, item_type, event_type, period_hour_unix) DO UPDATE SET total_count = total_count + 1`
	if _, err := singleton.D1Exec(sql, itemID, itemType, eventType, unixHour); err != nil {
		return err
	}
	return nil
}

func GetMetrics(itemID, itemType string) (*models.Metrics, error) {
	sql := `SELECT event_type, SUM(total_count) AS total_count
		FROM item_snapshots
		WHERE item_id = ? AND item_type = ?
		GROUP BY event_type`

	raw, err := singleton.D1Exec(sql, itemID, itemType)
	if err != nil {
		return nil, err
	}

	var rows []struct {
		EventType string `json:"event_type"`
		Total     int    `json:"total_count"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}

	metrics := &models.Metrics{
		ItemID:      itemID,
		ItemType:    itemType,
		ViewsCount:  0,
		LikesCount:  0,
		SharesCount: 0,
	}

	for _, row := range rows {
		switch row.EventType {
		case "view":
			metrics.ViewsCount = row.Total
		case "like":
			metrics.LikesCount = row.Total
		case "share":
			metrics.SharesCount = row.Total
		}
	}

	metrics.UpdatedAt = time.Now().Unix()
	return metrics, nil
}

func GetHistogram(itemID, itemType, eventType string, resolution string, from, to int64) ([]models.HistogramBucket, error) {
	secs, ok := models.ResolutionSeconds[resolution]
	if !ok {
		return nil, fmt.Errorf("invalid resolution %q: must be one of 1h, 1d, 1w, 1M", resolution)
	}
	if secs < models.ResolutionSeconds["1h"] {
		return nil, fmt.Errorf("resolution must be at least 1h")
	}
	if to <= 0 {
		to = time.Now().Unix()
	}
	if from <= 0 {
		from = to - 30*86400
	}
	if from >= to {
		return nil, fmt.Errorf("invalid time window: from must be less than to")
	}

	var bucketExpr string
	if resolution == "1M" {
		bucketExpr = "strftime('%s', strftime('%Y-%m-01', period_hour_unix, 'unixepoch'))"
	} else {
		bucketExpr = fmt.Sprintf("(period_hour_unix / %d) * %d", secs, secs)
	}

	sql := fmt.Sprintf(`SELECT %s AS bucket, SUM(total_count) AS total
		FROM item_snapshots
		WHERE item_id = ? AND item_type = ? AND event_type = ?
		AND period_hour_unix >= ? AND period_hour_unix < ?
		GROUP BY bucket
		ORDER BY bucket ASC`, bucketExpr)

	raw, err := singleton.D1Exec(sql, itemID, itemType, eventType, from, to)
	if err != nil {
		return nil, err
	}

	var rows []models.HistogramBucket
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}
