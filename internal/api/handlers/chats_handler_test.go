package handlers

import (
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

const baseURL = "http://localhost:50151/api/"

func TestChatCreate(t *testing.T) {

	var currentURL = baseURL + "chats"

	requestChat := &chat.Chat{
		ID:    config.ChatID1,
		Title: "Chat 24",
		Members: []chat.Member{
			{
				UserID:   config.Mahdi,
				IsActive: true,
				Role:     "reza",
				JoinedAt: time.Now(),
			},
			{
				UserID:   config.Golnar,
				IsActive: true,
				Role:     "creator",
				JoinedAt: time.Now(),
			},
			{
				UserID:   config.Ali,
				IsActive: true,
				Role:     "creator",
				JoinedAt: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "0",
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
		t.Errorf("Caption mismatch (-want +got):\n%s", diff)
	}
}

func TestChatRead(t *testing.T) {

	var currentURL = baseURL + "chats"

	queryParams := map[string]interface{}{
		"userId": config.Mahdi,
		"chatId": config.Varzesh3,
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

func TestChatReadAll(t *testing.T) {

	var currentURL = baseURL + "chats"

	queryParams := map[string]interface{}{
		"userId":    config.Mahdi,
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

func TestChatUpdate(t *testing.T) {

	var currentURL = baseURL + "chats"

	requestChat := &chat.Chat{
		ID:    config.ChatID1,
		Title: "Chat 24",
		Members: []chat.Member{
			{
				UserID:   config.Mahdi,
				IsActive: true,
				Role:     "reza",
				JoinedAt: time.Now(),
			},
			{
				UserID:   config.Golnar,
				IsActive: true,
				Role:     "creator",
				JoinedAt: time.Now(),
			},
			{
				UserID:   config.Ali,
				IsActive: true,
				Role:     "creator",
				JoinedAt: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "1",
	}

	respBody, err := helpers.MakeRequest(t, "PATCH", currentURL, nil, requestChat)
	if err != nil {
		t.Errorf("update request failed: %v", err)
	}

	var createdChat chat.Chat
	if err := json.Unmarshal(respBody, &createdChat); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	if diff := cmp.Diff(createdChat.Title, createdChat.Title); diff != "" {
		t.Errorf("Caption mismatch (-want +got):\n%s", diff)
	}
}

func TestChatDelete(t *testing.T) {

	var currentURL = baseURL + "chats/" + "0199357b-8352-7a0c-b168-4740fc60eb74"

	respBody, err := helpers.MakeRequest(t, "DELETE", currentURL, nil, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	t.Logf("Deleted %d chat", respBody)
}
