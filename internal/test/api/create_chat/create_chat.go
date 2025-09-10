// This program demonstrates how to create an HTTP client in Go
// to make a GET request to a specific endpoint with query parameters,
// using a struct to define the parameters.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

func createChat(newChat chat.Chat) error {

	fullURL := "http://localhost:50151/api/chats/"

	// 2. Marshal the struct into a JSON byte slice.
	// This converts your Go struct into raw data.
	jsonData, err := json.Marshal(newChat)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// 3. Create a new io.Reader from the JSON byte slice.
	// This makes the data streamable for the HTTP request.
	bodyReader := bytes.NewReader(jsonData)

	// 4. Create the new PATCH request.
	req, err := http.NewRequest("POST", fullURL, bodyReader)
	if err != nil {
		fmt.Println("Error creating PATCH request:", err)
		return err
	}

	// 5. Set the Content-Type header on the request object.
	req.Header.Set("Content-Type", "application/json")

	// 6. Use an http.Client to send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making PATCH request:", err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	fmt.Println("Status code:", resp.Status)

	return nil
}

func main() {

	members := []chat.Member{
		{
			UserID:     config.Mahdi,
			Role:       "creator",
			IsActive:   true,
			LastActive: time.Now(),
			JoinedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			UserID:     config.Parsa,
			Role:       "member",
			IsActive:   true,
			LastActive: time.Now(),
			JoinedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			UserID:     config.Golnar,
			Role:       "member",
			IsActive:   true,
			LastActive: time.Now(),
			JoinedAt:   time.Now().Add(-12 * time.Hour),
		},
		{
			UserID:     config.Behzad,
			Role:       "admin",
			IsActive:   false,
			LastActive: time.Now(),
			JoinedAt:   time.Now().Add(-12 * time.Hour),
		},
	}

	newChat := chat.Chat{
		Title:      "Chat  with Mahyar",
		IsVerified: true,
		Members:    members,
	}

	err := createChat(newChat)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//fmt.Println("Response received:")
}
