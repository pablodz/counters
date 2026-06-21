package counterscli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	BaseURL    = "http://localhost:8080"
	httpClient = &http.Client{Timeout: 5 * time.Second}
)

type Event struct {
	ItemID    string `json:"item_id"`
	ItemType  string `json:"item_type"`
	EventType string `json:"event_type"`
}

type Metrics struct {
	ItemID      string `json:"item_id"`
	ItemType    string `json:"item_type"`
	ViewsCount  int    `json:"views_count"`
	LikesCount  int    `json:"likes_count"`
	SharesCount int    `json:"shares_count"`
	UpdatedAt   int64  `json:"updated_at"`
}

type HistogramBucket struct {
	Bucket int64 `json:"bucket"`
	Total  int   `json:"total"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

func SetBaseURL(url string) {
	BaseURL = url
}

func IncrementEvent(event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	reqBody, err := http.NewRequest("POST", BaseURL+"/api/v1/metrics", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	reqBody.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		log.Printf("Increment event failed: %d %s", resp.StatusCode, string(bodyBytes))
		return &APIError{StatusCode: resp.StatusCode, Message: string(bodyBytes)}
	}
	return nil
}

func GetMetrics(itemType, itemID string) (*Metrics, error) {
	url := BaseURL + "/api/v1/" + itemType + "/" + itemID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		log.Printf("Get metrics failed: %d %s", resp.StatusCode, string(bodyBytes))
		return nil, &APIError{StatusCode: resp.StatusCode, Message: string(bodyBytes)}
	}

	var result Metrics
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetHistogram(itemType, itemID, eventType, resolution string, from, to int64) ([]HistogramBucket, error) {
	url := BaseURL + "/api/v1/histogram/" + itemType + "/" + itemID + "/" + eventType + "?resolution=" + resolution
	if from > 0 {
		url += fmt.Sprintf("&from=%d", from)
	}
	if to > 0 {
		url += fmt.Sprintf("&to=%d", to)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		log.Printf("Get histogram failed: %d %s", resp.StatusCode, string(bodyBytes))
		return nil, &APIError{StatusCode: resp.StatusCode, Message: string(bodyBytes)}
	}

	var result []HistogramBucket
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SetHTTPClient overrides the internal HTTP client (useful for testing).
func SetHTTPClient(client *http.Client) {
	httpClient = client
}
