package collection_manager_gemini

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// CollectionItem is the interface that every item in the collection must implement.
// The ID is of type uuid.UUID for type safety.
type CollectionItem interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
}

// Manager is the main struct for managing the collection.
type Manager[T CollectionItem] struct {
	baseDir string
	items   *Registry[T]
	mu      sync.RWMutex
}

// NewCollectionManager creates a new instance of Manager.
func NewCollectionManager[T CollectionItem](path string) (*Manager[T], error) {

	if strings.HasSuffix(path, ".json") {
		return nil, errors.New("path must be a directory, not a file")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	manager := &Manager[T]{
		baseDir: path,
		items:   NewRegistry[T](),
	}

	items, err := manager.readAllItemsFromDisk()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	for _, item := range items {
		manager.items.Register(item.GetID().String(), item)
	}

	return manager, nil
}

// itemPath returns the path to a JSON file based on the item's ID.
func (m *Manager[T]) itemPath(id uuid.UUID) string {
	return filepath.Join(m.baseDir, id.String()+".json")
}

// ---
// Helper functions that access the disk directly and are called by public methods.
// They do not have their own locks.

func (m *Manager[T]) readAllItemsFromDisk() ([]T, error) {
	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filename := strings.TrimSuffix(entry.Name(), ".json")
		if _, err := uuid.Parse(filename); err != nil {
			fmt.Printf("collection_manager: skipping file with invalid UUID filename: %s, error: %v\n", entry.Name(), err)
			continue
		}

		item, err := m.readItemFromDisk(filename)
		if err != nil {
			fmt.Printf("collection_manager: error reading item %s: %v\n", filename, err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (m *Manager[T]) readItemFromDisk(id string) (T, error) {
	var zero T
	path := filepath.Join(m.baseDir, id+".json")
	file, err := os.ReadFile(path)
	if err != nil {
		return zero, err
	}

	if len(file) == 0 {
		return zero, errors.New("empty file")
	}

	var item T
	if err := json.Unmarshal(file, &item); err != nil {
		return zero, err
	}
	return item, nil
}

func (m *Manager[T]) writeItemToDisk(item T) error {
	path := m.itemPath(item.GetID())
	jsonData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}

	tempFile := path + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return err
	}
	return os.Rename(tempFile, path)
}

// ---
// Public methods that manage concurrency and call helper functions.

// Create a new item in the collection.
func (m *Manager[T]) Create(newItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if reflect.ValueOf(newItem).IsNil() {
		var zero T
		return zero, errors.New("cannot create nil item")
	}

	if _, err := m.items.Get(newItem.GetID().String()); err == nil {
		var zero T
		return zero, fmt.Errorf("item with ID %s already exists", newItem.GetID().String())
	}

	if err := m.writeItemToDisk(newItem); err != nil {
		var zero T
		return zero, err
	}

	m.items.Register(newItem.GetID().String(), newItem)
	return newItem, nil
}

// Read an item from the collection by its ID.
func (m *Manager[T]) Read(id uuid.UUID) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.items.Get(id.String())
}

// ReadAll items from the collection.
func (m *Manager[T]) ReadAll() ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.items.GetAllValues(), nil
}

// Update an existing item in the collection.
func (m *Manager[T]) Update(updatedItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if reflect.ValueOf(updatedItem).IsNil() {
		var zero T
		return zero, errors.New("cannot update with nil item")
	}

	if _, err := m.items.Get(updatedItem.GetID().String()); err != nil {
		var zero T
		return zero, fmt.Errorf("item with ID %s does not exist", updatedItem.GetID().String())
	}

	if err := m.writeItemToDisk(updatedItem); err != nil {
		var zero T
		return zero, err
	}

	m.items.Update(updatedItem.GetID().String(), updatedItem)
	return updatedItem, nil
}

// Delete an item from the collection by its ID.
func (m *Manager[T]) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := m.items.Get(id.String()); err != nil {
		return fmt.Errorf("item with ID %s does not exist", id.String())
	}

	path := m.itemPath(id)
	if err := os.Remove(path); err != nil {
		return err
	}

	m.items.Delete(id.String())
	return nil
}

// ---
// Registry section: in-memory data store.
// This registry uses string keys, which is fine as uuid.UUID is always converted to string for storage.

type Registry[T any] struct {
	items map[string]T
	mu    sync.RWMutex
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

func (r *Registry[T]) Register(key string, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = value
}

func (r *Registry[T]) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
}

func (r *Registry[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]T)
}

func (r *Registry[T]) Get(key string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exists := r.items[key]
	if !exists {
		var zero T
		return zero, errors.New("key not found")
	}
	return val, nil
}

func (r *Registry[T]) Update(key string, newValue T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = newValue
}

func (r *Registry[T]) GetAllValues() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]T, 0, len(r.items))
	for _, v := range r.items {
		result = append(result, v)
	}
	return result
}

func (r *Registry[T]) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items) == 0
}
