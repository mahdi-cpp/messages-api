package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

func TestMessageCreate(t *testing.T) {

	var currentURL = baseURL + "messages"

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

	respBody, err := helpers.MakeRequest(t, "POST", currentURL, nil, testMessage)
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

	var currentURL = baseURL + "messages"

	queryParams := map[string]interface{}{
		"userId":    config.UserId,
		"chatId":    config.ChatID,
		"messageId": config.MessageID,
	}
	respBody, err := helpers.MakeRequest(t, "GET", currentURL, queryParams, nil)
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

	var currentURL = baseURL + "messages"

	queryParams := map[string]interface{}{
		"userId":    config.UserId,
		"chatId":    config.ChatID,
		"offset":    0,
		"limit":     100,
		"sortBy":    "id",
		"sortOrder": "start",
	}

	respBody, err := helpers.MakeRequest(t, "GET", currentURL, queryParams, nil)
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

	var currentURL = baseURL + "messages"
	var text = "Golnar Message 1010"

	testMessage := &message.UpdateOptions{
		MessageType: "message",
		UserID:      config.UserId,
		ChatID:      config.ChatID,
		MessageID:   "01992e25-4ba9-73ae-9f26-bdfd0d4bceb9",
		Content:     text,
	}

	respBody, err := helpers.MakeRequest(t, "PATCH", currentURL, nil, testMessage)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var updatedMessage message.Message
	if err := json.Unmarshal(respBody, &updatedMessage); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	if diff := cmp.Diff(updatedMessage.Content, text); diff != "" {
		t.Errorf("Content mismatch (-want +got):\n%s", diff)
	}

	duration := time.Since(start)
	t.Logf("TestMessageUpdate took: %v", duration)
}
