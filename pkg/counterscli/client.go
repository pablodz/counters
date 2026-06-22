package counterscli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pablodz/counters/data/models"
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
	url := fmt.Sprintf("%s/api/v1/%s/%s/%s/%s", BaseURL, itemType, itemID, eventType, userId)
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
		return fmt.Errorf("increment event failed: %d %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

func GetMetrics(itemType, itemID string) (models.Metrics, error) {
	url := fmt.Sprintf("%s/api/v1/%s/%s", BaseURL, itemType, itemID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.Metrics{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return models.Metrics{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Metrics{}, err
	}
	if resp.StatusCode >= 400 {
		return models.Metrics{}, fmt.Errorf("get metrics failed: %d %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.Metrics
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return models.Metrics{}, err
	}

	return result, nil
}

func GetHistogram(itemType, itemID, resolution string) ([]models.HistogramBucket, error) {
	url := fmt.Sprintf("%s/api/v1/histogram/%s/%s?resolution=%s", BaseURL, itemType, itemID, resolution)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("get histogram failed: %d %s", resp.StatusCode, string(bodyBytes))
	}

	var result []models.HistogramBucket
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetRecentActivity(itemType, itemID string) ([]models.AuditLogPayload, error) {
	url := fmt.Sprintf("%s/api/v1/activity/%s/%s", BaseURL, itemType, itemID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("get recent activity failed: %d %s", resp.StatusCode, string(bodyBytes))
	}

	var result []models.AuditLogPayload
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
