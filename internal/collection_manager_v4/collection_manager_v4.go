package collection_manager_v4

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CollectionItem interface {
	SetID(string)
	GetID() string
}

type storage[T CollectionItem] interface {
	ReadAll(requireExist bool) ([]T, error)
	CreateItem(item T) error
	UpdateItem(item T) error
	DeleteItem(id string) error
}

type singleFileStorage[T CollectionItem] struct {
	ctrl *Control[[]T]
}

func (s *singleFileStorage[T]) ReadAll(requireExist bool) ([]T, error) {
	dataPtr, err := s.ctrl.Read(requireExist)
	if err != nil {
		return nil, err
	}
	if dataPtr == nil {
		return []T{}, nil
	}
	return *dataPtr, nil
}

func (s *singleFileStorage[T]) CreateItem(item T) error {
	items, err := s.ReadAll(false)
	if err != nil {
		return err
	}
	items = append(items, item)
	return s.ctrl.Write(&items)
}

func (s *singleFileStorage[T]) UpdateItem(updatedItem T) error {
	items, err := s.ReadAll(false)
	if err != nil {
		return err
	}
	found := false
	for i, item := range items {
		if item.GetID() == updatedItem.GetID() {
			items[i] = updatedItem
			found = true
			break
		}
	}
	if !found {
		return errors.New("item not found")
	}
	return s.ctrl.Write(&items)
}

func (s *singleFileStorage[T]) DeleteItem(id string) error {
	items, err := s.ReadAll(false)
	if err != nil {
		return err
	}
	var newItems []T
	for _, item := range items {
		if item.GetID() != id {
			newItems = append(newItems, item)
		}
	}
	return s.ctrl.Write(&newItems)
}

type directoryStorage[T CollectionItem] struct {
	baseDir string
}

func (d *directoryStorage[T]) itemPath(id string) string {
	return filepath.Join(d.baseDir, id+".json")
}

func (d *directoryStorage[T]) ReadAll(requireExist bool) ([]T, error) {

	if _, err := os.Stat(d.baseDir); err != nil {
		if os.IsNotExist(err) {
			if requireExist {
				return nil, err
			}
			return []T{}, nil
		}
		return nil, err
	}

	entries, err := os.ReadDir(d.baseDir)
	if err != nil {
		return nil, err
	}

	var items []T
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if filepath.Ext(filename) != ".json" {
			continue
		}

		id := strings.TrimSuffix(filename, filepath.Ext(filename))
		item, err := d.readItem(id)
		if err != nil {
			fmt.Println("collection_manager_v4:", err.Error())
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *directoryStorage[T]) readItem(id string) (T, error) {
	var zero T
	path := d.itemPath(id)
	ctrl := NewMetadataManager[T](path)

	dataPtr, err := ctrl.Read(true)
	if err != nil {
		return zero, err
	}
	if dataPtr == nil {
		return zero, errors.New("metadata not found")
	}
	return *dataPtr, nil
}

func (d *directoryStorage[T]) CreateItem(item T) error {
	if err := os.MkdirAll(d.baseDir, 0755); err != nil {
		return err
	}
	path := d.itemPath(item.GetID())
	ctrl := NewMetadataManager[T](path)
	return ctrl.Write(&item)
}

func (d *directoryStorage[T]) UpdateItem(item T) error {
	path := d.itemPath(item.GetID())
	ctrl := NewMetadataManager[T](path)
	return ctrl.Write(&item)
}

func (d *directoryStorage[T]) DeleteItem(id string) error {
	path := d.itemPath(id)
	return os.Remove(path)
}

type Manager[T CollectionItem] struct {
	storage storage[T]
	items   *Registry[T]
}

func NewCollectionManager[T CollectionItem](path string, requireExist bool) (*Manager[T], error) {
	var store storage[T]

	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() {
			store = &directoryStorage[T]{baseDir: path}
		} else {
			store = &singleFileStorage[T]{ctrl: NewMetadataManager[[]T](path)}
		}
	} else {
		if strings.HasSuffix(path, ".json") {
			store = &singleFileStorage[T]{ctrl: NewMetadataManager[[]T](path)}
		} else {
			store = &directoryStorage[T]{baseDir: path}
		}
	}

	manager := &Manager[T]{
		storage: store,
		items:   NewRegistry[T](),
	}

	items, err := manager.storage.ReadAll(requireExist)
	if err != nil {
		return nil, fmt.Errorf("failed to load items: %w", err)
	}

	for _, item := range items {
		manager.items.Register(item.GetID(), item)
	}

	return manager, nil
}

func (manager *Manager[T]) Create(newItem T) (T, error) {

	if err := manager.storage.CreateItem(newItem); err != nil {
		return newItem, err
	}

	manager.items.Register(newItem.GetID(), newItem)
	return newItem, nil
}

func (manager *Manager[T]) Read(id string) (T, error) {
	return manager.items.Get(id)
}

func (manager *Manager[T]) ReadAll() ([]T, error) {
	return manager.items.GetAllValues(), nil
}

func (manager *Manager[T]) Update(updatedItem T) (T, error) {

	if err := manager.storage.UpdateItem(updatedItem); err != nil {
		return updatedItem, err
	}
	manager.items.Update(updatedItem.GetID(), updatedItem)
	return updatedItem, nil
}

func (manager *Manager[T]) Delete(id string) error {
	if err := manager.storage.DeleteItem(id); err != nil {
		return err
	}
	manager.items.Delete(id)
	return nil
}

type Control[T any] struct {
	filePath string
	mutex    sync.RWMutex
}

func NewMetadataManager[T any](filePath string) *Control[T] {
	return &Control[T]{
		filePath: filePath,
	}
}

func (control *Control[T]) Read(requireExist bool) (*T, error) {
	control.mutex.RLock()
	defer control.mutex.RUnlock()

	data := new(T)
	file, err := os.ReadFile(control.filePath)
	if err != nil {
		if os.IsNotExist(err) && requireExist {
			return nil, fmt.Errorf("file %s does not exist", control.filePath)
		}
		if os.IsNotExist(err) { // If not requiring existence, return empty data
			return data, nil
		}
		return nil, err
	}

	if len(file) == 0 {
		return data, nil
	}

	if err := json.Unmarshal(file, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (control *Control[T]) Update(updateFunc func(*T) error) error {
	control.mutex.Lock()
	defer control.mutex.Unlock()

	// Read current data into a pointer
	data, err := control.readData()
	if err != nil {
		return err
	}

	// Apply updates to the pointer
	if err := updateFunc(data); err != nil {
		return err
	}

	// Write updated data
	return control.writeData(data)
}

func (control *Control[T]) readData() (*T, error) {
	data := new(T)
	file, err := os.ReadFile(control.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}

	if len(file) > 0 {
		if err := json.Unmarshal(file, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (control *Control[T]) writeData(data *T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	tempFile := control.filePath + ".tmp"
	if err := os.WriteFile(tempFile, jsonData, 0644); err != nil {
		return err
	}

	return os.Rename(tempFile, control.filePath)
}

func (control *Control[T]) Write(data *T) error {
	control.mutex.Lock()
	defer control.mutex.Unlock()
	return control.writeData(data)
}

//Registry Section ----------------------------------------------------------

// Registry uses type parameters at struct level instead of method level
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
	r.mu.Lock()         // Acquire write lock (since we're modifying the map)
	defer r.mu.Unlock() // Ensure the lock is released

	//if _, exists := r.items[key]; !exists {
	//	return fmt.Errorf("key '%s' not found", key)
	//}

	delete(r.items, key) // Remove the key from the map
	//return nil
}

func (r *Registry[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = make(map[string]T) // Reinitialize the map
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
	r.mu.RLock()         // Acquire read lock for thread safety
	defer r.mu.RUnlock() // Ensure the lock is released

	return len(r.items) == 0
}
