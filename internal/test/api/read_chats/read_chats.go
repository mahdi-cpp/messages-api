// This program demonstrates how to create an HTTP client in Go
// to make a GET request to a specific endpoint with query parameters,
// using a struct to define the parameters.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
)

// buildQueryParams uses reflection to convert a struct into url.Values.
// It iterates through the struct's fields, checking for the `form` tag
// and skipping fields with `omitempty` if their value is the zero value for their type.
func buildQueryParams(s interface{}) (url.Values, error) {
	params := url.Values{}
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct or a pointer to a struct, got %s", v.Kind())
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		// Check for omitempty and skip if the value is zero
		if strings.HasSuffix(formTag, ",omitempty") {
			formTag = strings.TrimSuffix(formTag, ",omitempty")
			if value.IsZero() {
				continue
			}
		}

		// Handle different types to format them correctly for the URL.
		switch value.Kind() {
		case reflect.String:
			params.Add(formTag, value.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			params.Add(formTag, strconv.FormatInt(value.Int(), 10))
		case reflect.Bool:
			params.Add(formTag, strconv.FormatBool(value.Bool()))
		case reflect.Slice:
			// Special handling for slices. Check the type of the slice elements.
			elemType := value.Type().Elem()
			if elemType.Kind() == reflect.Struct {
				// We don't want to include a slice of structs in the query parameters.
				// It's a bad practice for GET requests. We'll simply skip this field.
				continue
			}

			// For simple slices (e.g., []string), add each element as a separate value.
			for j := 0; j < value.Len(); j++ {
				params.Add(formTag, fmt.Sprintf("%v", value.Index(j).Interface()))
			}
		case reflect.Ptr:
			if !value.IsNil() {
				// Dereference the pointer to get the actual value.
				switch value.Elem().Kind() {
				case reflect.Bool:
					params.Add(formTag, strconv.FormatBool(value.Elem().Bool()))
				case reflect.Slice:
					// Handle pointer to slice, check element type
					sliceValue := value.Elem()
					elemType := sliceValue.Type().Elem()
					if elemType.Kind() == reflect.Struct {
						continue // Skip slice of structs
					}
					for j := 0; j < sliceValue.Len(); j++ {
						params.Add(formTag, fmt.Sprintf("%v", sliceValue.Index(j).Interface()))
					}
				case reflect.Struct:
					// For struct types like time.Time.
					if _, ok := value.Interface().(*time.Time); ok {
						params.Add(formTag, value.Elem().Interface().(time.Time).Format(time.RFC3339))
					} else {
						// For other structs, we'll skip them as well.
						continue
					}
				}
			}
		}
	}
	return params, nil
}

// getChats makes a GET request to the chat API with the specified search options.
func getChats(options chat.SearchOptions) (string, error) {

	baseURL := "http://localhost:50151/api/chats"

	// Use the new helper function to build the query parameters automatically.
	params, err := buildQueryParams(options)
	if err != nil {
		return "", fmt.Errorf("failed to build query parameters: %w", err)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Println("Making request to URL:", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func main() {

	//isVerified := true

	options := chat.SearchOptions{
		//IsVerified: &isVerified,
		Offset:    0,
		Limit:     1,
		SortBy:    "id",
		SortOrder: "start",
	}

	response, err := getChats(options)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	//for _, chat1 := range response {
	//	fmt.Println(chat1.)
	//}

	fmt.Println("Response received:")
	fmt.Println(response)
}
