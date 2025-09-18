package handlers

import (
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

func TestMessageCreate(t *testing.T) {

	var currentURL = baseURL + "messages"

	testMessage := &message.Message{
		UserID:    config.Mahdi,
		ChatID:    config.ChatID1,
		Caption:   "Test Message 1001",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "1",
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
	if diff := cmp.Diff(testMessage.Caption, createdMessage.Caption); diff != "" {
		t.Errorf("Caption mismatch (-want +got):\n%s", diff)
	}
}

func TestMessageRead(t *testing.T) {

	var currentURL = baseURL + "messages"

	config.Init()

	queryParams := map[string]interface{}{
		"userId":    config.Mahdi,
		"chatId":    config.ChatID1,
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

	t.Logf("Retrieved %s", messages.Caption)
}

func TestMessageReadAll(t *testing.T) {

	var currentURL = baseURL + "messages"

	queryParams := map[string]interface{}{
		"userId":    config.Mahdi,
		"chatId":    config.ChatID1,
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
		UserID:    config.Mahdi,
		ChatID:    config.ChatID1,
		MessageID: config.MessageID,
		Content:   text,
	}

	respBody, err := helpers.MakeRequest(t, "PATCH", currentURL, nil, testMessage)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var updatedMessage message.Message
	if err := json.Unmarshal(respBody, &updatedMessage); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	if diff := cmp.Diff(updatedMessage.Caption, text); diff != "" {
		t.Errorf("Caption mismatch (-want +got):\n%s", diff)
	}

	duration := time.Since(start)
	t.Logf("TestMessageUpdate took: %v", duration)
}
