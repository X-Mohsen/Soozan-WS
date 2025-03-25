package interservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"Soozan-ws/config"
)

// HttpRequest sends a request to other service and returns response/error
func httpRequest(url string, method, token string, data interface{}) (*http.Response, error) {
	var requestBody []byte
	if data != nil {
		var err error
		requestBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	authValue := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", authValue)

	client := &http.Client{
		Timeout: 5 * time.Second, // TODO check timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func LoadChats(token string, page float64) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/request/inter-service/load-chats", config.APIBASEURL)
	method := "POST"

	data := map[string]interface{}{
		"page": page,
	}
	resp, err := httpRequest(endpointURL, method, token, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func LoadMessages(token string, requestID, page float64) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/request/inter-service/load-messages", config.APIBASEURL)
	method := "POST"

	data := map[string]interface{}{
		"page": page,
	}
	resp, err := httpRequest(endpointURL, method, token, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func NewMessage(token string, requestID float64, data map[string]interface{}) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/request/inter-service/new-message", config.APIBASEURL)
	method := "POST"

	resp, err := httpRequest(endpointURL, method, token, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func DeleteMessage(token string, messageID float64) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/request/inter-service/delete-message", config.APIBASEURL)
	method := "POST"

	var data []byte

	resp, err := httpRequest(endpointURL, method, token, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
