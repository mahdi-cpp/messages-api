package handlers

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

//func makeMessageRequest(t *testing.T, method, endpoint string, queryParams map[string]interface{}, body interface{}) ([]byte, error) {
//
//	// Build URL with query parameters
//	u, err := url.Parse(baseURL + endpoint)
//	if err != nil {
//		return nil, fmt.Errorf("parsing URL: %w", err)
//	}
//
//	if queryParams != nil {
//		q := u.Query()
//		for key, value := range queryParams {
//			q.Add(key, fmt.Sprintf("%v", value))
//		}
//		u.RawQuery = q.Encode()
//	}
//
//	// Marshal body if provided
//	var bodyReader io.Reader
//	if body != nil {
//		jsonData, err := json.Marshal(body)
//		if err != nil {
//			return nil, fmt.Errorf("marshaling body: %w", err)
//		}
//		bodyReader = bytes.NewReader(jsonData)
//	}
//
//	// Create request
//	req, err := http.NewRequest(method, u.String(), bodyReader)
//	if err != nil {
//		return nil, fmt.Errorf("creating request: %w", err)
//	}
//	if body != nil {
//		req.Header.Set("Content-Type", "application/json")
//	}
//
//	// Execute request
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, fmt.Errorf("executing request: %w", err)
//	}
//	defer resp.Body.Close()
//
//	// Check status code
//	if resp.StatusCode >= 400 {
//		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status)
//	}
//
//	// ReadChat response
//	respBody, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return nil, fmt.Errorf("reading response: %w", err)
//	}
//
//	return respBody, nil
//}

func TestMessageCreate(t *testing.T) {

	testMessage := &message.Message{
		MessageType: "message",
		Width:       450,
		UserID:      config.UserId,
		ChatID:      config.ChatID,
		Content:     "Test Message 1001",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     "1",
	}

	respBody, err := makeRequest(t, "POST", "messages", nil, testMessage)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var createdMessage message.Message
	if err := json.Unmarshal(respBody, &createdMessage); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	t.Logf("Created message ID: %s", createdMessage.ID)
	if diff := cmp.Diff(testMessage.Content, createdMessage.Content); diff != "" {
		t.Errorf("Content mismatch (-want +got):\n%s", diff)
	}
}

func TestMessageRead(t *testing.T) {

	queryParams := map[string]interface{}{
		"userId":    config.UserId,
		"chatId":    config.ChatID,
		"messageId": config.MessageID,
	}
	respBody, err := makeRequest(t, "GET", "messages", queryParams, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var messages message.Message
	if err := json.Unmarshal(respBody, &messages); err != nil {
		t.Fatalf("Unmarshaling failed: %v", err)
	}

	t.Logf("Retrieved %s", messages.Content)
}

func TestMessageReadAll(t *testing.T) {

	queryParams := map[string]interface{}{
		"userId":    config.UserId,
		"chatId":    config.ChatID,
		"offset":    0,
		"limit":     100,
		"sortBy":    "id",
		"sortOrder": "start",
	}

	respBody, err := makeRequest(t, "GET", "messages", queryParams, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var messages []message.Message
	if err := json.Unmarshal(respBody, &messages); err != nil {
		t.Fatalf("Unmarshaling failed: %v", err)
	}

	t.Logf("Retrieved %d messages", len(messages))
}

func TestMessageUpdate(t *testing.T) {

	start := time.Now()

	var text = "Ali Message 1009"

	testMessage := &message.UpdateOptions{
		MessageType: "message",
		UserID:      config.UserId,
		ChatID:      config.ChatID,
		MessageID:   "01992e25-4ba9-73ae-9f26-bdfd0d4bceb9",
		Content:     text,
	}

	respBody, err := makeRequest(t, "PATCH", "messages", nil, testMessage)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var updatedMessage message.Message
	if err := json.Unmarshal(respBody, &updatedMessage); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	fmt.Println("message Content: ", updatedMessage.Content)

	//if diff := cmp.Diff(updatedMessage.Content, text); diff != "" {
	//	t.Errorf("Content mismatch (-want +got):\n%s", diff)
	//}

	duration := time.Since(start)
	t.Logf("TestMessageUpdate took: %v", duration)
}
