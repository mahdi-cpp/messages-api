package main

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

func createMessage(newChat *message.Message) (*message.Message, error) {

	fullURL := "http://localhost:50151/api/messages/"

	// 2. Marshal the struct into a JSON byte slice.
	// This converts your Go struct into raw data.
	jsonData, err := json.Marshal(newChat)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	// 3. create a new io.Reader from the JSON byte slice.
	// This makes the data streamable for the HTTP request.
	bodyReader := bytes.NewReader(jsonData)

	// 4. create the new POST request.
	req, err := http.NewRequest("POST", fullURL, bodyReader)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// 5. Set the Caption-Type header on the request object.
	req.Header.Set("Content-Type", "application/json")

	// 6. Use an http.Client to send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}

	// 7. Use a defer statement to ensure the response body is closed.
	// This prevents resource leaks.
	defer resp.Body.Close()

	// 8. Handle non-2xx status codes by returning an error.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received error status code: %s", resp.Status)
	}

	// 9. ReadChat the response body into a byte slice.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 10. Unmarshal the JSON response body into a new message.Message struct.
	var createdMessage message.Message
	err = json.Unmarshal(body, &createdMessage)
	if err != nil {
		fmt.Println("Error unmarshaling response JSON:", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	fmt.Println("Status code:", resp.Status)

	// 11. Return the unmarshaled struct and a nil error.
	return &createdMessage, nil
}

func TestCreate(t *testing.T) {

	sendMessage := &message.Message{
		MessageType: "message",
		Width:       450,
		UserID:      config.Mahdi,
		ChatID:      config.ChatID1,
		Caption:     "Test Message",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Version:     "1",
	}

	getMessage, err := createMessage(sendMessage)
	if err != nil {
		t.Fatalf("Error creating message: %v", err)
		return
	}

	// 3. Use cmp.Equal to perform a deep comparison of all fields.
	if !cmp.Equal(sendMessage.Caption, getMessage.Caption) {
		// If the structs are not equal, use cmp.Diff to get a human-readable
		// diff showing exactly which fields differ.
		t.Errorf("GetDefaultUser() returned an unexpected message.\nDiff:\n%s", cmp.Diff(sendMessage, getMessage))
		return
	}

	fmt.Println("Created message:", getMessage.ID)

}
