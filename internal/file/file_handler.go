package collection_manager_db

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

// Set record size to 2 kilobytes
const (
	recordSize    = 2048 // 2 KB
	dataFileName  = "data.db"
	indexFileName = "index.db"
)

// FileHandler is a struct for managing resources
type FileHandler struct {
	file  *os.File
	index map[uuid.UUID]int64 // Map UUID to the data file offset
}

// NewFileHandler opens/creates the data and index files
func NewFileHandler() (*FileHandler, error) {
	file, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening data file: %w", err)
	}

	handler := &FileHandler{
		file:  file,
		index: make(map[uuid.UUID]int64),
	}

	// Load the index from disk
	handler.loadIndex()
	return handler, nil
}

// Close closes the files and saves the index
func (h *FileHandler) Close() error {
	// Save the index before closing
	if err := h.saveIndex(); err != nil {
		return fmt.Errorf("error saving index: %w", err)
	}
	return h.file.Close()
}

// saveIndex saves the current index to a JSON file
func (h *FileHandler) saveIndex() error {
	data, err := json.Marshal(h.index)
	if err != nil {
		return fmt.Errorf("error serializing index: %w", err)
	}
	return os.WriteFile(indexFileName, data, 0644)
}

// loadIndex loads the index from the JSON file
func (h *FileHandler) loadIndex() {
	data, err := os.ReadFile(indexFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Index file does not exist, creating a new index.")
			return
		}
		log.Printf("error reading index file: %v\n", err)
		return
	}
	if err := json.Unmarshal(data, &h.index); err != nil {
		log.Printf("error deserializing index: %v\n", err)
	}
}

// WriteRecord writes a new record and returns a UUID
func (h *FileHandler) WriteRecord(data []byte) (uuid.UUID, error) {
	if len(data) > recordSize {
		return uuid.Nil, fmt.Errorf("data size is larger than max record size (2KB)")
	}

	// Seek to the end of the file to write the new record
	offset, err := h.file.Seek(0, io.SeekEnd)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error seeking to end of file: %w", err)
	}

	// Create a Version 7 UUID
	u, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creating Version 7 UUID: %w", err)
	}

	// Store the offset in the index
	h.index[u] = offset

	// Create a 2KB buffer
	recordBuffer := make([]byte, recordSize)
	copy(recordBuffer, data)
	// If the data is smaller than 2KB, the rest of the buffer is padded with zeros.

	// Write the full buffer to the data file
	if _, err := h.file.Write(recordBuffer); err != nil {
		return uuid.Nil, fmt.Errorf("error writing record: %w", err)
	}

	return u, nil
}

// ReadRecord reads a record by its UUID
func (h *FileHandler) ReadRecord(recordUUID uuid.UUID) ([]byte, error) {
	offset, ok := h.index[recordUUID]
	if !ok {
		return nil, fmt.Errorf("UUID not found: %s", recordUUID)
	}

	recordBuffer := make([]byte, recordSize)
	n, err := h.file.ReadAt(recordBuffer, offset)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading block from offset %d: %w", offset, err)
	}

	// To ensure that we return only the actual data,
	// a mechanism is needed to specify the actual data length,
	// such as a length header at the beginning of the record or
	// searching for a null byte. The entire 2KB buffer is currently returned.
	return recordBuffer[:n], nil
}

// ReadAllRecords reads all records from the data file
func (h *FileHandler) ReadAllRecords() (map[uuid.UUID][]byte, error) {
	// 1. Extract UUID keys from the index map
	uuids := make([]uuid.UUID, 0, len(h.index))
	for u := range h.index {
		uuids = append(uuids, u)
	}

	// 2. Sort the keys by offset to read them in the order they were written
	sort.Slice(uuids, func(i, j int) bool {
		return h.index[uuids[i]] < h.index[uuids[j]]
	})

	records := make(map[uuid.UUID][]byte)
	for _, u := range uuids {
		record, err := h.ReadRecord(u)
		if err != nil {
			// If an error occurs while reading a record, log it and continue
			log.Printf("error reading record with UUID %s: %v\n", u.String(), err)
			continue
		}
		records[u] = record
	}

	return records, nil
}
