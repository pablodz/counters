package counterscli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pablodz/counters/data/models"
	"github.com/pablodz/counters/handlers"
)

var (
	BaseURL    = "http://localhost:8080"
	httpClient = &http.Client{Timeout: 5 * time.Second}
)

func SetBaseURL(url string) {
	BaseURL = url
}

func IncrementLike(itemId, itemType, userId string) error {
	return incrementEvent(itemType, itemId, "like", userId)
}

func IncrementUnlike(itemId, itemType, userId string) error {
	return incrementEvent(itemType, itemId, "unlike", userId)
}

func IncrementShare(itemId, itemType, userId string) error {
	return incrementEvent(itemType, itemId, "share", userId)
}

func IncrementView(itemId, itemType, userId string) error {
	return incrementEvent(itemType, itemId, "view", userId)
}

func incrementEvent(itemType, itemID, eventType, userId string) error {
	url := BaseURL + fmt.Sprintf(handlers.IncrementEventURL, itemType, itemID, eventType, userId)
	return get(url, nil)
}

func get(url string, result any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed: %d %s", resp.StatusCode, string(bodyBytes))
	}

	return json.Unmarshal(bodyBytes, result)
}

func GetMetrics(itemType, itemID string) (models.Metrics, error) {
	url := BaseURL + fmt.Sprintf(handlers.GetMetricsURL, itemType, itemID)
	var result models.Metrics
	err := get(url, &result)
	return result, err
}

func GetHistogram(itemType, itemID, resolution string) ([]models.HistogramBucket, error) {
	url := BaseURL + fmt.Sprintf(handlers.GetHistogramURL, itemType, itemID) + "?resolution=" + resolution
	var result []models.HistogramBucket
	err := get(url, &result)
	return result, err
}

func GetMetricsList(itemType string, itemIDs []string) (map[string]models.Metrics, error) {
	idsStr := strings.Join(itemIDs, ",")
	url := BaseURL + fmt.Sprintf(handlers.GetMetricsListURL, itemType, idsStr)
	var result map[string]models.Metrics
	err := get(url, &result)
	return result, err
}

func GetRecentActivity(itemType, itemID string) ([]models.AuditLogPayload, error) {
	url := BaseURL + fmt.Sprintf(handlers.GetRecentActivityURL, itemType, itemID)
	var result []models.AuditLogPayload
	err := get(url, &result)
	return result, err
}
