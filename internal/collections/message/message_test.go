package message

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collection_manager"
	"github.com/mahdi-cpp/messages-api/internal/collections/metadata"
	"github.com/mahdi-cpp/messages-api/internal/config"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

const baseURL = "http://localhost:50103/api/v1/upload/"

// httpClient is a shared instance of the HTTP client for efficiency.
var httpClient = &http.Client{Timeout: 30 * time.Second}

type Response struct {
	Message  string    `json:"message,omitempty"`
	Filename string    `json:"filename,omitempty"`
	ID       uuid.UUID `json:"id,omitempty"`
	Error    string    `json:"error,omitempty"`
}

type DirectoryRequest struct {
	ID      uuid.UUID `json:"id"`
	Message string    `json:"message"`
	Errors  string    `json:"errors,omitempty"`
}

type Request struct {
	Directory uuid.UUID `json:"directory"`
	IsVideo   bool      `json:"isVideo"`
	//Hash      string    `json:"hash"`
}

func TestReadMessages(t *testing.T) {

	var err error
	var messagesDirectory = config.GetPath("/test/messages")
	collectionManager, err := collection_manager.New[*Message](messagesDirectory)
	if err != nil {
		t.Error(err)
	}

	all, err := collectionManager.ReadAll()
	if err != nil {
		return
	}

	for _, msg := range all {
		if msg.Medias != nil {
			fmt.Println(msg.Caption, msg.Medias[0].MimeType)
		}
	}
}

var uploadsDir = "/app/iris/services/uploads"
var appDir = "/app/iris/com.iris.messages/chats"

func TestMessageCreate(t *testing.T) {

	config.Init()

	chatID := config.ChatID3
	var messagesDirectory = filepath.Join(config.RootDir, "chats", chatID.String(), "metadata/v1/messages")

	var err error

	collectionManager, err := collection_manager.New[*Message](messagesDirectory)
	if err != nil {
		t.Error(err)
	}

	var apiURL = baseURL + "create"
	respBody, err := helpers.MakeRequest(t, "POST", apiURL, nil, nil)
	if err != nil {
		t.Errorf("create request failed: %v", err)
	}

	var workDir DirectoryRequest
	if err := json.Unmarshal(respBody, &workDir); err != nil {
		t.Errorf("unmarshaling response: %v", err)
	}

	fmt.Println(workDir.ID)
	apiURL = baseURL + "media"

	// A context with a 30-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//files := []string{
	//	"/app/tmp/f/0.jpg",
	//	"/app/tmp/f/1.jpg",
	//	"/app/tmp/f/2.jpg",
	//	"/app/tmp/f/3.jpg",
	//	"/app/tmp/f/4.mp4",
	//	"/app/tmp/f/5.mp4",
	//}

	files := []string{
		"/app/tmp/huma/1.jpg",
		"/app/tmp/huma/2.jpg",
		"/app/tmp/huma/3.jpg",
		"/app/tmp/huma/4.jpg",
		"/app/tmp/huma/5.mp4",
	}

	var medias []*Media
	for i := 0; i < len(files); i++ {
		mediaMetadata, err := upload(ctx, httpClient, apiURL, workDir.ID, files[i])
		if err != nil {
			t.Errorf("%v", err)
		}
		if mediaMetadata == nil || mediaMetadata.ID == uuid.Nil {
			t.Fatal("Expected a non-nil response, but got nil")
		}
		medias = append(medias, mediaMetadata)
	}

	for _, media := range medias {
		moveMedia(chatID.String(), workDir.ID.String(), media)
	}

	msg := &Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		UserID:    config.Mahdi,
		Caption:   "Huma Group",
		Directory: filepath.Join(config.RootDir, "chats", chatID.String()),
		AssetType: "media",
		Medias:    medias,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "1",
	}

	create, err := collectionManager.Create(msg)
	if err != nil {
		t.Errorf("Error creating message: %s", err)
	}

	fmt.Println("message created: ", create.ID)
}

func upload(ctx context.Context, client *http.Client, apiURL string, directoryID uuid.UUID, filePath string) (*Media, error) {

	var mimeType = ""
	if isJPEG(filePath) {
		mimeType = "image/jpeg"
	} else {
		mimeType = "video/mp4"
	}

	// Create the Request struct with the necessary metadata.
	var uploadReq = &Request{}
	if mimeType == "image/jpeg" {
		uploadReq = &Request{
			Directory: directoryID,
			IsVideo:   false,
		}
	} else {
		uploadReq = &Request{
			Directory: directoryID,
			IsVideo:   true,
		}
	}

	// Call the mediaUpload function.
	resp, err := mediaUpload(ctx, client, apiURL, filePath, uploadReq)
	if err != nil {
		return nil, fmt.Errorf("media upload failed: %w", err)
	}
	var media = &Media{}
	if mimeType == "image/jpeg" {
		media = &Media{
			ID:          resp.ID,
			FileSize:    resp.FileInfo.FileSize,
			MimeType:    resp.FileInfo.MimeType,
			Width:       resp.Image.Width,
			Height:      resp.Image.Height,
			Orientation: resp.Image.Orientation,
		}
	} else {
		media = &Media{
			ID:          resp.ID,
			FileSize:    resp.FileInfo.FileSize,
			MimeType:    resp.FileInfo.MimeType,
			Width:       resp.Video.Width,
			Height:      resp.Video.Height,
			Orientation: "",
		}
	}

	return media, nil
}

func mediaUpload(ctx context.Context, client *http.Client, apiURL, filePath string, uploadRequest *Request) (*metadata.Metadata, error) {

	// Open the file to be uploaded.
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Create a new multipart writer.
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Marshal the struct into a JSON string.
	jsonData, err := json.Marshal(uploadRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Create a new form field for the JSON data.
	if err := writer.WriteField("metadata", string(jsonData)); err != nil {
		return nil, fmt.Errorf("failed to write JSON data to form field: %w", err)
	}

	// Create a form file part for the media.
	filePart, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy the file content into the form file part.
	if _, err := io.Copy(filePart, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the multipart writer to finalize the body.
	writer.Close()

	// Create the HTTP request.
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Caption-Type header with the boundary from the writer.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for a successful response status code.
	if resp.StatusCode != http.StatusOK {
		// Read and log the server's error message.
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("server responded with status code %d, but could not read error body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("server responded with status code %d and body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the JSON response body.
	var serverResponse metadata.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&serverResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	fmt.Printf("Successfully uploaded image from %s\n", filePath)
	return &serverResponse, nil
}

func isJPEG(filePath string) bool {
	// Convert the file path to lowercase for a case-insensitive check.
	lowerCasePath := strings.ToLower(filePath)

	return strings.HasSuffix(lowerCasePath, ".jpg") || strings.HasSuffix(lowerCasePath, ".jpeg")
}

func moveMedia(chatID string, workDir string, media *Media) {

	var id = media.ID.String()
	if media.MimeType == "image/jpeg" {
		var name = id + ".jpg"
		src := filepath.Join(uploadsDir, workDir, name)
		des := filepath.Join(config.RootDir, "chats", chatID, "assets", name)
		err := os.Rename(src, des)
		if err != nil {
			return
		}

		var thumbnail200 = id + "-200.jpg"
		srcThumb := filepath.Join(uploadsDir, workDir, thumbnail200)
		desThumb := filepath.Join(config.RootDir, "chats", chatID, "thumbnails", thumbnail200)
		err = os.Rename(srcThumb, desThumb)
		if err != nil {
			return
		}
	} else {
		var name = id + ".mp4"
		src := filepath.Join(uploadsDir, workDir, name)
		des := filepath.Join(config.RootDir, "chats", chatID, "assets", name)
		err := os.Rename(src, des)
		if err != nil {
			return
		}

		var thumbnail200 = id + "-200.jpg"
		srcThumb := filepath.Join(uploadsDir, workDir, thumbnail200)
		desThumb := filepath.Join(config.RootDir, "chats", chatID, "thumbnails", thumbnail200)
		err = os.Rename(srcThumb, desThumb)
		if err != nil {
			return
		}

		var thumbnail400 = id + "-400.jpg"
		srcThumb = filepath.Join(uploadsDir, workDir, thumbnail400)
		desThumb = filepath.Join(config.RootDir, "chats", chatID, "thumbnails", thumbnail400)
		err = os.Rename(srcThumb, desThumb)
		if err != nil {
			return
		}
	}

}
