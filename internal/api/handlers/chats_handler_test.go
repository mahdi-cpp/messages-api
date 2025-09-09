package handlers

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
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
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

	fmt.Println(u.String())
	fmt.Println("")

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

func createChat(t *testing.T, newChat *chat.Chat) (*chat.Chat, error) {
	respBody, err := makeRequest(t, "POST", "messages", nil, newChat)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	var createdChat chat.Chat
	if err := json.Unmarshal(respBody, &createdChat); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	t.Logf("Created message ID: %s", createdChat.ID)
	return &createdChat, nil
}

func TestCreate(t *testing.T) {

	testChat := &chat.Chat{
		ID:        config.ChatID,
		Title:     "Test Chat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}

	createdChat, err := createChat(t, testChat)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}

	if diff := cmp.Diff(testChat.Title, createdChat.Title); diff != "" {
		t.Errorf("Content mismatch (-want +got):\n%s", diff)
	}
}

func TestRead(t *testing.T) {

	queryParams := map[string]interface{}{
		"userId": config.UserId,
		"chatId": "018f3a8b-1b32-729a-f7e5-5467c1b2d3e4",
	}

	respBody, err := makeRequest(t, "GET", "chats", queryParams, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var chat1 chat.Chat
	if err := json.Unmarshal(respBody, &chat1); err != nil {
		t.Fatalf("Unmarshaling failed: %v", err)
	}

	t.Logf("Retrieved chat: %v", chat1)

}

func TestReadAll(t *testing.T) {

	queryParams := map[string]interface{}{
		"userId":    config.UserId,
		"offset":    2,
		"limit":     100,
		"sortBy":    "id",
		"sortOrder": "start",
	}

	respBody, err := makeRequest(t, "GET", "chats", queryParams, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var chats []chat.Chat
	if err := json.Unmarshal(respBody, &chats); err != nil {
		t.Fatalf("Unmarshaling failed: %v", err)
	}

	t.Logf("Retrieved %d chats", len(chats))
}
