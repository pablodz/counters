package singleton

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type d1RequestBody struct {
	SQL    string `json:"sql"`
	Params []any  `json:"params,omitempty"`
}

type d1Result struct {
	Success bool            `json:"success"`
	Results json.RawMessage `json:"results"`
	Error   string          `json:"error"`
}

type d1ResponseBody struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result []d1Result `json:"result"`
}

var d1Client = &http.Client{Timeout: 15 * time.Second}

func D1Exec(sql string, params ...any) (json.RawMessage, error) {
	if CF_ACCOUNT_ID == "" || CF_D1_DATABASE_ID == "" {
		return nil, fmt.Errorf("missing D1 endpoint configuration")
	}
	if D1_API_TOKEN == "" {
		return nil, fmt.Errorf("missing D1 API token")
	}

	apiURL := fmt.Sprintf(
		"https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/query",
		CF_ACCOUNT_ID, CF_D1_DATABASE_ID,
	)

	jsonData, err := json.Marshal(d1RequestBody{SQL: sql, Params: params})
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+D1_API_TOKEN)

	resp, err := d1Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("D1 HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var responseBody d1ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if !responseBody.Success {
		if len(responseBody.Errors) > 0 {
			return nil, fmt.Errorf("D1 error: %s", responseBody.Errors[0].Message)
		}
		return nil, fmt.Errorf("D1 request failed")
	}

	if len(responseBody.Result) == 0 {
		return json.RawMessage("[]"), nil
	}

	var failures []string
	for i, result := range responseBody.Result {
		if !result.Success {
			errMsg := result.Error
			if errMsg == "" {
				errMsg = "unknown error"
			}
			failures = append(failures, fmt.Sprintf("statement %d: %s", i, errMsg))
		}
	}
	if len(failures) > 0 {
		return nil, fmt.Errorf("query execution failed:\n%s", strings.Join(failures, "\n"))
	}

	if len(responseBody.Result[0].Results) == 0 || string(responseBody.Result[0].Results) == "null" {
		return json.RawMessage("[]"), nil
	}

	return responseBody.Result[0].Results, nil
}
