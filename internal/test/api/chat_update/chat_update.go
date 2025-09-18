// This program demonstrates how to create an HTTP client in Go
// to make a GET request to a specific endpoint with query parameters,
// using a struct to define the parameters.

package main

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"log"
	"net/http"

	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
)

func updateChats(chatID string, options chat.UpdateOptions) error {

	fullURL := "http://localhost:50151/api/chats/" + chatID

	// 2. Marshal the struct into a JSON byte slice.
	// This converts your Go struct into raw data.
	jsonData, err := json.Marshal(options)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// 3. create a new io.Reader from the JSON byte slice.
	// This makes the data streamable for the HTTP request.
	bodyReader := bytes.NewReader(jsonData)

	// 4. create the new PATCH request.
	req, err := http.NewRequest("PATCH", fullURL, bodyReader)
	if err != nil {
		fmt.Println("Error creating PATCH request:", err)
		return err
	}

	// 5. Set the Caption-Type header on the request object.
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

	//isVerified := true

	//members := []chat.Member{
	//	{
	//		UserID:     "20",
	//		Role:       "creator",
	//		IsActive:   true,
	//		LastActive: time.Now(),
	//		JoinedAt:   time.Now().Add(-24 * time.Hour),
	//	},
	//	{
	//		UserID:     "21",
	//		Role:       "member",
	//		IsActive:   true,
	//		LastActive: time.Now(),
	//		JoinedAt:   time.Now().Add(-12 * time.Hour),
	//	},
	//	{
	//		UserID:     "22",
	//		Role:       "member",
	//		IsActive:   true,
	//		LastActive: time.Now(),
	//		JoinedAt:   time.Now().Add(-12 * time.Hour),
	//	},
	//	{
	//		UserID:     "23",
	//		Role:       "admin",
	//		IsActive:   false,
	//		LastActive: time.Now(),
	//		JoinedAt:   time.Now().Add(-12 * time.Hour),
	//	},
	//}

	options := chat.UpdateOptions{
		Title: "Ali Tesla Team",
		//IsVerified: &isVerified,
		//AddMembers: members,
	}

	err := updateChats("018f3a8b-1b32-7293-c1d4-8765f4d1e2f3", options)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//fmt.Println("Response received:")
}
