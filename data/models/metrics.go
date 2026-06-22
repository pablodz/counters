package models

import "slices"

type HistogramBucket struct {
	Bucket    int64  `json:"bucket"`
	EventType string `json:"event_type"`
	Total     int64  `json:"total"`
}

var ValidResolutions = []string{"1h", "1d"}

func IsValidResolution(resolution string) bool {
	return slices.Contains(ValidResolutions, resolution)
}

type AuditLogPayload struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	UserType  string `json:"user_type"`
	ItemID    string `json:"item_id"`
	ItemType  string `json:"item_type"`
	EventType string `json:"event_type"`
	CreatedAt int64  `json:"created_at"`
}
