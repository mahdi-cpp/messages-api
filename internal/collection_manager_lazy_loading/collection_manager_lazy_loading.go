package collection_manager_lazy_loading

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	recordStatusSize = 1
	recordSize       = 2048 // 2 KB
	dataFileName     = "/app/tmp/messages/data.db"
	indexFileName    = "/app/tmp/messages/index.db"
	compactInterval  = 1 * time.Hour
)

const (
	StatusActive  = 0x00
	StatusDeleted = 0x01
)

// collectionItem is the interface that every item in the collection must implement.
type collectionItem interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
}

type collectionIndex interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
}

// FileHandler is a struct for managing file resources and the index.
type FileHandler struct {
	file  *os.File
	index map[uuid.UUID]int64
	mu    sync.RWMutex
}

func NewFileHandler() (*FileHandler, error) {
	file, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening data file: %w", err)
	}

	handler := &FileHandler{
		file:  file,
		index: make(map[uuid.UUID]int64),
	}

	handler.loadIndex()
	return handler, nil
}

func (h *FileHandler) Close() error {
	if err := h.saveIndex(); err != nil {
		log.Printf("Error saving index on close: %v", err)
	}
	return h.file.Close()
}

func (h *FileHandler) saveIndex() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stringIndex := make(map[string]int64)
	for k, v := range h.index {
		stringIndex[k.String()] = v
	}

	data, err := json.Marshal(stringIndex)
	if err != nil {
		return fmt.Errorf("error serializing index: %w", err)
	}

	tempFile := indexFileName + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("error writing temporary index file: %w", err)
	}

	return os.Rename(tempFile, indexFileName)
}

func (h *FileHandler) loadIndex() {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(indexFileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Index file does not exist. A new one will be created.")
			return
		}
		log.Printf("Error reading index file, attempting to rebuild: %v", err)
		h.rebuildIndex()
		return
	}

	var stringIndex map[string]int64
	if err := json.Unmarshal(data, &stringIndex); err != nil {
		log.Printf("Error deserializing index, attempting to rebuild: %v", err)
		h.rebuildIndex()
		return
	}

	h.index = make(map[uuid.UUID]int64, len(stringIndex))
	for k, v := range stringIndex {
		if id, err := uuid.Parse(k); err == nil {
			h.index[id] = v
		} else {
			log.Printf("Error parsing UUID %s from index: %v", k, err)
		}
	}

	log.Printf("Loaded %d entries from index.", len(h.index))
}

// WriteRecord writes a new record and returns a UUID.
func (h *FileHandler) WriteRecord(data []byte, id uuid.UUID) (uuid.UUID, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(data) > recordSize-recordStatusSize {
		return uuid.Nil, fmt.Errorf("data size is larger than max record size (%d bytes)", recordSize-recordStatusSize)
	}

	offset, err := h.file.Seek(0, io.SeekEnd)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error seeking to end of file: %w", err)
	}

	recordBuffer := make([]byte, recordSize)
	recordBuffer[0] = StatusActive
	copy(recordBuffer[recordStatusSize:], data)

	if _, err := h.file.Write(recordBuffer); err != nil {
		return uuid.Nil, fmt.Errorf("error writing record: %w", err)
	}

	h.index[id] = offset
	return id, nil
}

func (h *FileHandler) ReadRecord(recordUUID uuid.UUID) ([]byte, error) {
	h.mu.RLock()
	offset, ok := h.index[recordUUID]
	h.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("UUID not found: %s", recordUUID)
	}

	recordBuffer := make([]byte, recordSize)
	if _, err := h.file.ReadAt(recordBuffer, offset); err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading block from offset %d: %w", offset, err)
	}

	if recordBuffer[0] == StatusDeleted {
		return nil, fmt.Errorf("record with UUID %s is marked as deleted", recordUUID)
	}

	dataLength := bytes.IndexByte(recordBuffer[recordStatusSize:], 0)
	if dataLength == -1 {
		dataLength = recordSize - recordStatusSize
	}

	return recordBuffer[recordStatusSize : recordStatusSize+dataLength], nil
}

func (h *FileHandler) HasRecord(recordUUID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, exists := h.index[recordUUID]
	return exists
}

func (h *FileHandler) UpdateRecord(recordUUID uuid.UUID, data []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	offset, ok := h.index[recordUUID]
	if !ok {
		return fmt.Errorf("UUID not found: %s", recordUUID)
	}

	if len(data) > recordSize-recordStatusSize {
		return fmt.Errorf("data size is larger than max record size (%d bytes)", recordSize-recordStatusSize)
	}

	recordBuffer := make([]byte, recordSize)
	recordBuffer[0] = StatusActive
	copy(recordBuffer[recordStatusSize:], data)

	if _, err := h.file.WriteAt(recordBuffer, offset); err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	return nil
}

func (h *FileHandler) DeleteRecord(recordUUID uuid.UUID) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	offset, ok := h.index[recordUUID]
	if !ok {
		return fmt.Errorf("UUID not found: %s", recordUUID)
	}

	if _, err := h.file.WriteAt([]byte{StatusDeleted}, offset); err != nil {
		return fmt.Errorf("error marking record as deleted: %w", err)
	}

	delete(h.index, recordUUID)
	return nil
}

func (h *FileHandler) GetAllUUIDs() []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	uuids := make([]uuid.UUID, 0, len(h.index))
	for u := range h.index {
		uuids = append(uuids, u)
	}
	return uuids
}

func (h *FileHandler) rebuildIndex() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.index = make(map[uuid.UUID]int64)
	fileInfo, err := h.file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}
	fileSize := fileInfo.Size()

	for offset := int64(0); offset < fileSize; offset += recordSize {
		recordBuffer := make([]byte, recordSize)
		n, err := h.file.ReadAt(recordBuffer, offset)
		if err != nil && err != io.EOF {
			log.Printf("Error reading record at offset %d: %v", offset, err)
			continue
		}

		if n == 0 || recordBuffer[0] == StatusDeleted {
			continue
		}

		dataLength := bytes.IndexByte(recordBuffer[recordStatusSize:], 0)
		if dataLength == -1 {
			dataLength = recordSize - recordStatusSize
		}
		data := recordBuffer[recordStatusSize : recordStatusSize+dataLength]

		var item struct {
			ID uuid.UUID `json:"id"`
		}
		if err := json.Unmarshal(data, &item); err != nil {
			log.Printf("Error unmarshaling record at offset %d: %v", offset, err)
			continue
		}

		if item.ID != uuid.Nil {
			h.index[item.ID] = offset
		}
	}

	log.Printf("Rebuilt index with %d entries.", len(h.index))
	return nil
}

func (h *FileHandler) Compact() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	tempFileName := dataFileName + ".tmp"
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		return fmt.Errorf("error creating temporary file: %w", err)
	}
	defer tempFile.Close()

	oldIndex := h.index
	newIndex := make(map[uuid.UUID]int64)

	var currentOffset int64 = 0
	for u, oldOffset := range oldIndex {
		recordBuffer := make([]byte, recordSize)
		_, err := h.file.ReadAt(recordBuffer, oldOffset)
		if err != nil {
			log.Printf("Error reading record %s for compaction: %v", u, err)
			continue
		}

		if _, err := tempFile.Write(recordBuffer); err != nil {
			log.Printf("Error writing record %s to temporary file: %v", u, err)
			continue
		}

		newIndex[u] = currentOffset
		currentOffset += recordSize
	}

	h.file.Close()
	if err := os.Rename(tempFileName, dataFileName); err != nil {
		return fmt.Errorf("error renaming temporary file: %w", err)
	}

	h.file, err = os.OpenFile(dataFileName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error reopening data file: %w", err)
	}

	h.index = newIndex
	if err := h.saveIndex(); err != nil {
		log.Printf("Error saving index after compaction: %v", err)
	}

	log.Printf("Data file compacted. %d active records remain.", len(h.index))
	return nil
}

// ---

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

func (r *registry[T]) read(key uuid.UUID) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, exists := r.items[key]
	return val, exists
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

func (r *registry[T]) count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}

// ---

// Manager is the main struct for managing the collection.
type Manager[T collectionItem] struct {
	fh    *FileHandler
	items *registry[T]
	mu    sync.RWMutex
}

func New[T collectionItem]() (*Manager[T], error) {
	fh, err := NewFileHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to create file handler: %w", err)
	}

	manager := &Manager[T]{
		fh:    fh,
		items: newRegistry[T](),
	}

	go manager.startCompactionRoutine()
	return manager, nil
}

func (m *Manager[T]) startCompactionRoutine() {
	ticker := time.NewTicker(compactInterval)
	defer ticker.Stop()
	for range ticker.C {
		log.Println("Starting compaction routine...")
		if err := m.fh.Compact(); err != nil {
			log.Printf("Compaction failed: %v", err)
		} else {
			log.Println("Compaction routine finished.")
		}
	}
}

func (m *Manager[T]) Close() error {
	return m.fh.Close()
}

func (m *Manager[T]) loadItemFromDisk(id uuid.UUID) (T, error) {
	var zero T

	data, err := m.fh.ReadRecord(id)
	if err != nil {
		return zero, fmt.Errorf("error reading record from disk: %w", err)
	}

	var item T
	if err := json.Unmarshal(data, &item); err != nil {
		return zero, fmt.Errorf("error unmarshaling record: %w", err)
	}

	// This is the crucial fix: set the ID on the unmarshaled item
	item.SetID(id)

	m.items.create(id, item)

	return item, nil
}

func (m *Manager[T]) Create(newItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Step 1: Generate the UUID
	id, err := uuid.NewV7()
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error generating UUID: %w", err)
	}

	// Step 2: Set the UUID on the new item
	newItem.SetID(id)

	// Step 3: Marshal the item (now with a UUID) to JSON
	data, err := json.Marshal(newItem)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error marshaling item: %w", err)
	}

	// Step 4: Write the data to the file
	_, err = m.fh.WriteRecord(data, id)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error writing record: %w", err)
	}

	// Step 5: Add to the in-memory cache
	m.items.create(id, newItem)
	return newItem, nil
}

func (m *Manager[T]) Read(id uuid.UUID) (T, error) {
	m.mu.RLock()
	item, exists := m.items.read(id)
	m.mu.RUnlock()

	if exists {
		return item, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists = m.items.read(id)
	if exists {
		return item, nil
	}

	return m.loadItemFromDisk(id)
}

func (m *Manager[T]) ReadAll() ([]T, error) {
	m.mu.RLock()
	uuids := m.fh.GetAllUUIDs()
	m.mu.RUnlock()

	items := make([]T, 0, len(uuids))
	var errs []error

	for _, id := range uuids {
		if item, exists := m.items.read(id); exists {
			items = append(items, item)
			continue
		}

		m.mu.Lock()
		item, err := m.loadItemFromDisk(id)
		m.mu.Unlock()
		if err != nil {
			log.Printf("Error loading item %s: %v", id, err)
			errs = append(errs, err)
			continue
		}
		items = append(items, item)
	}

	if len(errs) > 0 {
		return items, fmt.Errorf("some records could not be read: %v", errs)
	}

	return items, nil
}

func (m *Manager[T]) Update(updatedItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := updatedItem.GetID()

	if !m.fh.HasRecord(id) {
		var zero T
		return zero, fmt.Errorf("item with ID %s does not exist", id.String())
	}

	data, err := json.Marshal(updatedItem)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error marshaling item: %w", err)
	}

	if err := m.fh.UpdateRecord(id, data); err != nil {
		var zero T
		return zero, fmt.Errorf("error updating record: %w", err)
	}

	m.items.update(id, updatedItem)
	return updatedItem, nil
}

func (m *Manager[T]) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.fh.DeleteRecord(id); err != nil {
		return err
	}

	m.items.delete(id)
	return nil
}

func (m *Manager[T]) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.fh.index)
}
