package collection_manager_db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/google/uuid"
)

const (
	recordSize    = 2048 // 2 KB
	dataFileName  = "data.db"
	indexFileName = "index.db"
)

// FileHandler is a struct for managing resources
type FileHandler struct {
	file  *os.File
	index map[uuid.UUID]int64 // Map UUID to the data file offset
	mu    sync.RWMutex
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
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(h.index)
	if err != nil {
		return fmt.Errorf("error serializing index: %w", err)
	}
	return os.WriteFile(indexFileName, data, 0644)
}

// loadIndex loads the index from the JSON file
func (h *FileHandler) loadIndex() {
	h.mu.Lock()
	defer h.mu.Unlock()

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
	h.mu.Lock()
	defer h.mu.Unlock()

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
// ReadRecord reads a record by its UUID
func (h *FileHandler) ReadRecord(recordUUID uuid.UUID) ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	offset, ok := h.index[recordUUID]
	if !ok {
		return nil, fmt.Errorf("UUID not found: %s", recordUUID)
	}

	recordBuffer := make([]byte, recordSize)
	_, err := h.file.ReadAt(recordBuffer, offset)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading block from offset %d: %w", offset, err)
	}

	// Find the first null byte to determine the actual data length
	// This assumes the data doesn't contain null bytes, which is true for JSON
	dataLength := bytes.IndexByte(recordBuffer, 0)
	if dataLength == -1 {
		dataLength = recordSize // No null byte found, use full buffer
	}

	return recordBuffer[:dataLength], nil
}

// ReadAllRecords reads all records from the data file
func (h *FileHandler) ReadAllRecords() (map[uuid.UUID][]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

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

// collectionItem is the interface that every item in the collection must implement.
// The ID is of type uuid.UUID for type safety.
type collectionItem interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
}

// ---
// registry section: in-memory data store.
// This registry uses string keys, which is fine as uuid.UUID is always converted to string for storage.

type registry[T any] struct {
	items map[uuid.UUID]T
	mu    sync.RWMutex
}

func newRegistry[T any]() *registry[T] {
	return &registry[T]{items: make(map[uuid.UUID]T)}
}

func (r *registry[T]) create(key uuid.UUID, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = value
}

func (r *registry[T]) read(key uuid.UUID) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exists := r.items[key]
	if !exists {
		var zero T
		return zero, errors.New("key not found")
	}
	return val, nil
}

func (r *registry[T]) readAll() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]T, 0, len(r.items))
	for _, v := range r.items {
		result = append(result, v)
	}
	return result
}

func (r *registry[T]) update(key uuid.UUID, newValue T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = newValue
}

func (r *registry[T]) delete(key uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
}

// Count returns the number of items in the registry.
func (r *registry[T]) count() int {
	// Acquire a read lock to safely read the map.
	r.mu.RLock()
	// The defer statement ensures to unlock happens before the function returns.
	defer r.mu.RUnlock()
	return len(r.items)
}

func (r *registry[T]) clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[uuid.UUID]T)
}

func (r *registry[T]) isEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items) == 0
}

// ---

// Manager is the main struct for managing the collection.
type Manager[T collectionItem] struct {
	fh    *FileHandler
	items *registry[T]
	mu    sync.RWMutex
}

// New creates a new instance of Manager.
func New[T collectionItem]() (*Manager[T], error) {

	fh, err := NewFileHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to create file handler: %w", err)
	}

	manager := &Manager[T]{
		fh:    fh,
		items: newRegistry[T](),
	}

	//Load all items from disk
	records, err := fh.ReadAllRecords()
	if err != nil {
		return nil, fmt.Errorf("failed to read records: %w", err)
	}

	for id, data := range records {
		var item T
		if err := json.Unmarshal(data, &item); err != nil {
			log.Printf("error unmarshaling record %s: %v", id.String(), err)
			continue
		}
		manager.items.create(id, item)
	}

	return manager, nil
}

// Close closes the underlying file handler
func (m *Manager[T]) Close() error {
	return m.fh.Close()
}

// Create a new item in the collection.
func (m *Manager[T]) Create(newItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Marshal the item to JSON
	data, err := json.Marshal(newItem)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error marshaling item: %w", err)
	}

	// Write to file
	id, err := m.fh.WriteRecord(data)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error writing record: %w", err)
	}

	// Set the ID and store in memory
	newItem.SetID(id)
	m.items.create(id, newItem)

	return newItem, nil
}

// Read an item from the collection by its ID.
func (m *Manager[T]) Read(id uuid.UUID) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.items.read(id)
}

// ReadAll items from the collection.
func (m *Manager[T]) ReadAll() ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.items.readAll(), nil
}

// Update an existing item in the collection.
func (m *Manager[T]) Update(updatedItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := updatedItem.GetID()

	// Marshal the updated item to JSON
	data, err := json.Marshal(updatedItem)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error marshaling item: %w", err)
	}

	// Write to file (this will create a new record)
	newID, err := m.fh.WriteRecord(data)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error writing record: %w", err)
	}

	// Update the ID if it changed
	if newID != id {
		updatedItem.SetID(newID)
		m.items.delete(id)
	}

	m.items.create(newID, updatedItem)
	return updatedItem, nil
}

// Delete an item from the collection by its ID.
func (m *Manager[T]) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Note: This only removes from memory, not from the data file
	// To fully remove from disk, we would need to implement a deletion mechanism
	// in the FileHandler that can handle record deletion and compaction
	m.items.delete(id)
	return nil
}

// Count returns the number of items in the collection.
func (m *Manager[T]) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.items.count()
}
