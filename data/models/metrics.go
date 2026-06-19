package models

type Metrics struct {
	ContentID   string `json:"content_id"`
	ContentType string `json:"content_type"`
	ViewsCount  int    `json:"views_count"`
	LikesCount  int    `json:"likes_count"`
	SharesCount int    `json:"shares_count"`
	UpdatedAt   string `json:"updated_at"`
}

var allowedFields = map[string]bool{
	"views_count":  true,
	"likes_count":  true,
	"shares_count": true,
}

func ValidField(f string) bool {
	return allowedFields[f]
}
