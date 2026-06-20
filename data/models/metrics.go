package models

import "time"

type Metrics struct {
	ItemID      string `json:"item_id"`
	ItemType    string `json:"item_type"`
	ViewsCount  int    `json:"views_count"`
	LikesCount  int    `json:"likes_count"`
	SharesCount int    `json:"shares_count"`
	UpdatedAt   int64  `json:"updated_at"`
}

type TrackingPayload struct {
	ItemID    string `json:"item_id"`
	ItemType  string `json:"item_type"`
	EventType string `json:"event_type"`
}

type HistogramBucket struct {
	Bucket int64 `json:"bucket"`
	Total  int   `json:"total"`
}

// ResolutionSeconds maps a resolution token to its duration in seconds.
// The minimum supported resolution is 1 hour.
var ResolutionSeconds = map[string]int64{
	"1h": 3600,
	"1d": 86400,
	"1w": 604800,
	"1M": 2592000,
}

func PrepararDatosInteraccion(payload TrackingPayload) (int64, error) {
	return time.Now().Truncate(time.Hour).Unix(), nil
}
