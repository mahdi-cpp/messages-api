package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

const baseURL = "http://localhost:50151/api/"

// Helper function to make HTTP requests
func makeRequest(t *testing.T, method, endpoint string, queryParams map[string]interface{}, body interface{}) ([]byte, error) {

	// Build URL with query parameters
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	if queryParams != nil {
		q := u.Query()
		for key, value := range queryParams {
			q.Add(key, fmt.Sprintf("%v", value))
		}
		u.RawQuery = q.Encode()
	}

	// Marshal body if provided
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// Create request
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status)
	}

	// ReadChat response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return respBody, nil
}

func createMessage(t *testing.T, newMessage *Message) (*Message, error) {
	respBody, err := makeRequest(t, "POST", "messages", nil, newMessage)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	var createdMessage Message
	if err := json.Unmarshal(respBody, &createdMessage); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	t.Logf("Created message ID: %s", createdMessage.ID)
	return &createdMessage, nil
}

func TestCreate(t *testing.T) {

	testMessage := &Message{
		MessageType: "message",
		Width:       450,
		UserID:      config.Mahdi,
		ChatID:      config.TestChatID,
		Content:     "Test Message",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     "1",
	}

	createdMessage, err := createMessage(t, testMessage)
	if err != nil {
		t.Fatalf("Error creating message: %v", err)
	}

	if diff := cmp.Diff(testMessage.Content, createdMessage.Content); diff != "" {
		t.Errorf("Content mismatch (-want +got):\n%s", diff)
	}
}

func testReadMessages(t *testing.T, queryParams map[string]interface{}) {
	respBody, err := makeRequest(t, "GET", "messages", queryParams, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var messages []Message
	if err := json.Unmarshal(respBody, &messages); err != nil {
		t.Fatalf("Unmarshaling failed: %v", err)
	}

	//t.Logf(queryParams)
	t.Logf("Retrieved %d messages", len(messages))
}

func TestRead(t *testing.T) {
	queryParams := map[string]interface{}{
		"userId":    config.Mahdi,
		"chatId":    config.TestChatID,
		"content":   "\u001b[A8888",
		"offset":    0,
		"limit":     100,
		"sortBy":    "id",
		"sortOrder": "start",
	}
	testReadMessages(t, queryParams)
}

func TestReadAll(t *testing.T) {
	queryParams := map[string]interface{}{
		"userId":    config.Mahdi,
		"chatId":    config.TestChatID,
		"offset":    0,
		"limit":     100,
		"sortBy":    "id",
		"sortOrder": "start",
	}
	testReadMessages(t, queryParams)
}
