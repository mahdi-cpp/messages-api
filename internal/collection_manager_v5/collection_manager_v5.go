package collection_manager_v5

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type CollectionItem interface {
	SetID(string)
	GetID() string
}

type Manager[T CollectionItem] struct {
	baseDir string
	items   *registry[T]
	mu      sync.RWMutex
}

func New[T CollectionItem](path string) (*Manager[T], error) {

	// Ensure the path is a directory
	if strings.HasSuffix(path, ".json") {
		return nil, errors.New("path must be a directory, not a file")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	manager := &Manager[T]{
		baseDir: path,
		items:   newRegistry[T](),
	}

	// Load existing items
	items, err := manager.readAllItems()
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	// Create all loaded items
	for _, item := range items {
		manager.items.Create(item.GetID(), item)
	}

	return manager, nil
}

func (m *Manager[T]) itemPath(id string) string {
	return filepath.Join(m.baseDir, id+".json")
}

func (m *Manager[T]) readAllItems() ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, err
	}

	items := make([]T, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) != ".json" {
			continue
		}

		id := strings.TrimSuffix(filename, filepath.Ext(filename))
		item, err := m.readItem(id)
		if err != nil {
			fmt.Printf("collection_manager: error reading item %s: %v\n", id, err)
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (m *Manager[T]) readItem(id string) (T, error) {
	var zero T
	path := m.itemPath(id)

	m.mu.RLock()
	defer m.mu.RUnlock()

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

func (m *Manager[T]) writeItem(item T) error {
	path := m.itemPath(item.GetID())

	m.mu.Lock()
	defer m.mu.Unlock()

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

func (m *Manager[T]) Create(newItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if newItem == nil {
		return newItem, errors.New("can not create item, it is nil")
	}

	if !isValidUUID(newItem.GetID()) {
		return newItem, fmt.Errorf("invalid UUID format")
	}

	// Check if item already exists
	if _, err := m.items.Read(newItem.GetID()); err == nil {
		var zero T
		return zero, fmt.Errorf("item with ID %s already exists", newItem.GetID())
	}

	if err := m.writeItem(newItem); err != nil {
		var zero T
		return zero, err
	}

	m.items.Create(newItem.GetID(), newItem)
	return newItem, nil
}

func (m *Manager[T]) Read(id string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.items.Read(id)
}

func (m *Manager[T]) ReadAll() ([]T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.items.ReadAll(), nil
}

func (m *Manager[T]) Update(updatedItem T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if item exists
	if _, err := m.items.Read(updatedItem.GetID()); err != nil {
		var zero T
		return zero, fmt.Errorf("item with ID %s does not exist", updatedItem.GetID())
	}

	if err := m.writeItem(updatedItem); err != nil {
		var zero T
		return zero, err
	}

	m.items.Update(updatedItem.GetID(), updatedItem)
	return updatedItem, nil
}

func (m *Manager[T]) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if item exists
	if _, err := m.items.Read(id); err != nil {
		return fmt.Errorf("item with ID %s does not exist", id)
	}

	path := m.itemPath(id)
	if err := os.Remove(path); err != nil {
		return err
	}

	m.items.Delete(id)
	return nil
}

// registry Section ----------------------------------------------------------

type registry[T any] struct {
	items map[string]T
	mu    sync.RWMutex
}

func newRegistry[T any]() *registry[T] {
	return &registry[T]{items: make(map[string]T)}
}

func (r *registry[T]) Create(key string, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = value
}

func (r *registry[T]) Read(key string) (T, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, exists := r.items[key]
	if !exists {
		var zero T
		return zero, errors.New("key not found")
	}
	return val, nil
}

func (r *registry[T]) ReadAll() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]T, 0, len(r.items))
	for _, v := range r.items {
		result = append(result, v)
	}
	return result
}

func (r *registry[T]) Update(key string, newValue T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = newValue
}

func (r *registry[T]) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
}

func (r *registry[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = make(map[string]T)
}

func (r *registry[T]) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items) == 0
}

//utils-------------------------

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return true
}
