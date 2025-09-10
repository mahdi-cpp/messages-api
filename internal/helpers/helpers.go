package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserID از Gin context، user_id را به صورت string دریافت می‌کند.
func GetUserID(c *gin.Context) (string, bool) {

	//// این تابع باید بعد از middleware احراز هویت استفاده شود
	//userID, exists := c.Read("X-User-ID")
	//if !exists {
	//	return "", false
	//}
	//
	//userIDStr, ok := userID.(string)
	//if !ok {
	//	return "", false
	//}
	//
	//return userIDStr, true

	//// این تابع باید بعد از middleware احراز هویت استفاده شود
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return "", false
	}

	return userID, true
}

// MakeRequest Helper function to make HTTP requests
func MakeRequest(t *testing.T, method, endpoint string, queryParams map[string]interface{}, body interface{}) ([]byte, error) {

	// Build URL with query parameters
	u, err := url.Parse(endpoint)
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

// IsValidUUID checks if a given string is a valid UUID (v4 or v7).
// It returns a boolean and an error if parsing fails.
//func IsValidUUID(s uuid.UUID) (bool, error) {
//	_, err := uuid.Parse(s)
//	if err != nil {
//		return false, fmt.Errorf("invalid UUID format: %w", err)
//	}
//	return true, nil
//}

func GetFileSize(filepath string) (int64, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func GetUserId(c *gin.Context) (int, error) {

	userIDStr := c.GetHeader("userID")
	fmt.Println(userIDStr)

	return strconv.Atoi(userIDStr)
}

func GenerateUUID() (uuid.UUID, error) {
	u7, err2 := uuid.NewV7()
	if err2 != nil {
		return uuid.UUID{}, fmt.Errorf("error generating UUIDv7: %w", err2)
	}
	return u7, nil
}
