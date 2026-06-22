package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/singleton"
)

func LogInteraction(log models.AuditLogPayload) error {
	now := time.Now().Unix()
	unixHour := (now / 3600) * 3600

	// 1. Actualizar el contador por hora (Upsert)
	sqlHourly := `INSERT INTO item_interactions_hourly (item_id, item_type, event_type, period_hour_unix, total_count)
		VALUES (?, ?, ?, ?, 1)
		ON CONFLICT(item_id, item_type, event_type, period_hour_unix) DO UPDATE SET total_count = total_count + 1`

	if _, err := singleton.D1Exec(sqlHourly, log.ItemID, log.ItemType, log.EventType, unixHour); err != nil {
		return err
	}

	sqlAudit := `INSERT INTO user_item_interactions_log
		(user_id, user_type, item_id, item_type, event_type, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	if _, err := singleton.D1Exec(sqlAudit, log.UserID, log.UserType, log.ItemID, log.ItemType, log.EventType, log.CreatedAt); err != nil {
		return err
	}

	return nil
}

func GetMetrics(itemID, itemType string) (map[string]int, error) {
	sql := `SELECT event_type, SUM(total_count) AS total_count
		FROM item_interactions_hourly
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

	metrics := make(map[string]int)
	for _, row := range rows {
		metrics[row.EventType] = row.Total
	}

	return metrics, nil
}

func GetHistogram(itemID, itemType string, resolution string) ([]models.HistogramBucket, error) {
	to := time.Now().Unix()
	var points int64
	var secs int64

	switch resolution {
	case "1h":
		points = 24
		secs = 3600
	case "1d":
		points = 30
		secs = 86400
	default:
		points = 24
		secs = 3600
	}

	from := to - (points * secs)
	startBucket := (from / secs) * secs

	eventTypes := []string{"view", "like", "share"}
	buckets := make([]models.HistogramBucket, 0, points*int64(len(eventTypes)))
	bucketIndex := make(map[string]int)

	idx := 0
	for i := int64(0); i < points; i++ {
		bTime := startBucket + (i * secs)
		for _, ev := range eventTypes {
			buckets = append(buckets, models.HistogramBucket{
				Bucket:    bTime,
				EventType: ev,
				Total:     0,
			})
			key := fmt.Sprintf("%d_%s", bTime, ev)
			bucketIndex[key] = idx
			idx++
		}
	}

	bucketExpr := fmt.Sprintf("(period_hour_unix / %d) * %d", secs, secs)

	sql := fmt.Sprintf(`SELECT %s AS bucket, event_type, SUM(total_count) AS total
		FROM item_interactions_hourly
		WHERE item_id = ? AND item_type = ?
		AND period_hour_unix >= ? AND period_hour_unix < ?
		GROUP BY bucket, event_type
		ORDER BY bucket ASC`, bucketExpr)

	raw, err := singleton.D1Exec(sql, itemID, itemType, from, to)
	if err != nil {
		return buckets, err
	}

	var dbRows []models.HistogramBucket
	if err := json.Unmarshal(raw, &dbRows); err != nil {
		return buckets, err
	}

	for _, row := range dbRows {
		key := fmt.Sprintf("%d_%s", row.Bucket, row.EventType)
		if index, exists := bucketIndex[key]; exists {
			buckets[index].Total = row.Total
		}
	}

	return buckets, nil
}
