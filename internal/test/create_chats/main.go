package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
)

// ProcessUserJSONFile reads a JSON file containing an array of User structs,
// and for each user, creates a separate JSON file in the specified output directory.
// Each new file is named after the user's ID.
// Note: This function remains as a general utility for file-based processing.
func ProcessUserJSONFile(inputFilePath, outputDirPath string) error {

	// Read the entire input JSON file.
	// io/os.ReadFile is deprecated in Go 1.16+, replaced by os.ReadFile.
	// For wider compatibility or older Go versions, ioutil.ReadFile is acceptable.
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to read input file %s: %w", inputFilePath, err)
	}

	// Unmarshal the JSON data into a slice of User structs.
	var users []chat.Chat
	err = json.Unmarshal(data, &users)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", inputFilePath, err)
	}

	// Ensure the output directory exists. If it doesn't, create it.
	// The 0755 permission mode grants read/write/execute for the owner,
	// and read/execute for group and others.
	err = os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDirPath, err)
	}

	// Iterate through each selectUser in the slice.
	for _, selectUser := range users {
		// Marshal the single selectUser struct back into JSON, with indentation for readability.
		userJSON, err := json.MarshalIndent(selectUser, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal selectUser %s to JSON: %w", selectUser.ID, err)
		}

		// Construct the output file path using the selectUser's ID as the filename.
		outputFileName := fmt.Sprintf("%s.json", selectUser.ID)
		outputFilePath := filepath.Join(outputDirPath, outputFileName)

		// Write the marshaled selectUser JSON to the new file.
		// os.WriteFile is the preferred way to write data to a file in Go 1.16+.
		err = os.WriteFile(outputFilePath, userJSON, 0644) // 0644 for read/write for owner, read-only for others
		if err != nil {
			return fmt.Errorf("failed to write selectUser file %s: %w", outputFilePath, err)
		}

		fmt.Printf("Successfully created file: %s\n", outputFilePath)
	}

	return nil // No error occurred
}

func main() {
	// --- Example Usage: Directly processing JSON data from a string literal ---

	// Define the JSON data as a string. This replaces the dummyJSONData variable
	// and the need to write/read it from a file.

	filePath := "./internal/documents/chats.json"

	// Attempt to read the entire file
	content, err := os.ReadFile(filePath) // Reads the entire file into a byte slice
	if err != nil {
		// Handle potential errors like file not found or permission issues
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	outputDirName := "./internal/test/tmp/chats"

	// Unmarshal the JSON data directly from the string into a slice of User structs.
	var users []chat.Chat
	err = json.Unmarshal([]byte(content), &users)
	if err != nil {
		fmt.Printf("Error unmarshaling raw JSON data: %v\n", err)
		return
	}

	// Ensure the output directory exists.
	err = os.MkdirAll(outputDirName, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory %s: %v\n", outputDirName, err)
		return
	}

	// Iterate through each selectUser and create a separate JSON file.
	for _, selectUser := range users {
		userJSON, err := json.MarshalIndent(selectUser, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling selectUser %s to JSON: %v\n", selectUser.ID, err)
			continue
		}

		outputFileName := fmt.Sprintf("%s.json", selectUser.ID)
		outputFilePath := filepath.Join(outputDirName, outputFileName)

		err = os.WriteFile(outputFilePath, userJSON, 0644)
		if err != nil {
			fmt.Printf("Error writing selectUser file %s: %v\n", outputFilePath, err)
			continue
		}

		fmt.Printf("Successfully created file: %s\n", outputFilePath)
	}

	fmt.Println("User files created successfully from in-memory JSON data!")

	// Clean up: optionally remove the output directory
	// os.RemoveAll(outputDirName)
}
