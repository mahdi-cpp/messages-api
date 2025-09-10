package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

const baseURL = "http://localhost:50151/api/"

func TestChatCreate(t *testing.T) {

	var currentURL = baseURL + "chats"

	requestChat := &chat.Chat{
		ID:    config.ChatID,
		Title: "Chat 24",
		Members: []chat.Member{
			{
				UserID:   config.UserId,
				IsActive: true,
				Role:     "reza",
				JoinedAt: time.Now(),
			},
			{
				UserID:   "018f3a8b-1b32-729b-8f90-1234a5b6c7d8",
				IsActive: true,
				Role:     "creator",
				JoinedAt: time.Now(),
			},
			{
				UserID:   "018f3a8b-1b32-729c-a1b2-9876a5b4c3d2",
				IsActive: true,
				Role:     "creator",
				JoinedAt: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}

	respBody, err := helpers.MakeRequest(t, "POST", currentURL, nil, requestChat)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var createdChat chat.Chat
	if err := json.Unmarshal(respBody, &createdChat); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	if diff := cmp.Diff(createdChat.Title, createdChat.Title); diff != "" {
		t.Errorf("Content mismatch (-want +got):\n%s", diff)
	}
}

func TestRead(t *testing.T) {

	var currentURL = baseURL + "chats"

	queryParams := map[string]interface{}{
		"userId": config.UserId,
		"chatId": "018f3a8b-1b32-729a-f7e5-5467c1b2d3e4",
	}

	respBody, err := helpers.MakeRequest(t, "GET", currentURL, queryParams, nil)
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

	var currentURL = baseURL + "chats"

	queryParams := map[string]interface{}{
		"userId":    config.UserId,
		"offset":    2,
		"limit":     100,
		"sortBy":    "id",
		"sortOrder": "start",
	}

	respBody, err := helpers.MakeRequest(t, "GET", currentURL, queryParams, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	var chats []chat.Chat
	if err := json.Unmarshal(respBody, &chats); err != nil {
		t.Fatalf("Unmarshaling failed: %v", err)
	}

	t.Logf("Retrieved %d chats", len(chats))
}
