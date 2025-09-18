package collection_manager_generic_index

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/google/uuid"
)

const (
	recordStatusSize = 1
	recordSize       = 2048
	indexRecordSize  = 400
	dirName          = "/app/tmp/messages"
)

const (
	StatusActive  = 0x00
	StatusDeleted = 0x01
)

type collectionItem interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
}

type IndexEntry[I any] struct {
	Offset      int64 // موقعیت داده در فایل داده
	IndexOffset int64 // موقعیت رکورد ایندکس در فایل ایندکس
	IndexData   I
}

type FileHandler[I any] struct {
	dataFile  *os.File
	indexFile *os.File
	mu        sync.RWMutex
	dataPath  string
	indexPath string
}

func NewFileHandler[I any]() (*FileHandler[I], error) {

	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating directory %s: %w", dirName, err)
	}

	dataFileName := filepath.Join(dirName, "data.db")
	indexFileName := filepath.Join(dirName, "index.db")

	dataFile, err := os.OpenFile(dataFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening data file: %w", err)
	}

	indexFile, err := os.OpenFile(indexFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		dataFile.Close()
		return nil, fmt.Errorf("error opening index file: %w", err)
	}

	return &FileHandler[I]{
		dataFile:  dataFile,
		indexFile: indexFile,
		dataPath:  dataFileName,
		indexPath: indexFileName,
	}, nil
}

func (h *FileHandler[I]) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var errs []error
	if err := h.dataFile.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing data file: %w", err))
	}
	if err := h.indexFile.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing index file: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errs)
	}

	return nil
}

func (h *FileHandler[I]) WriteRecord(data []byte) (int64, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(data) > recordSize-recordStatusSize {
		return -1, fmt.Errorf("data size is larger than max record size (%d bytes)", recordSize-recordStatusSize)
	}

	offset, err := h.dataFile.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("error seeking to end of data file: %w", err)
	}

	recordBuffer := make([]byte, recordSize)
	recordBuffer[0] = StatusActive
	copy(recordBuffer[recordStatusSize:], data)

	if _, err := h.dataFile.Write(recordBuffer); err != nil {
		return -1, fmt.Errorf("error writing record: %w", err)
	}

	return offset, nil
}

func (h *FileHandler[I]) ReadRecord(offset int64) ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if offset < 0 {
		return nil, fmt.Errorf("invalid offset: %d", offset)
	}

	recordBuffer := make([]byte, recordSize)
	n, err := h.dataFile.ReadAt(recordBuffer, offset)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error reading block from data file at offset %d: %w", offset, err)
	}

	if n == 0 {
		return nil, fmt.Errorf("no data read at offset %d", offset)
	}

	if recordBuffer[0] == StatusDeleted {
		return nil, fmt.Errorf("record at offset %d is marked as deleted", offset)
	}

	dataLength := bytes.IndexByte(recordBuffer[recordStatusSize:], 0)
	if dataLength == -1 {
		dataLength = recordSize - recordStatusSize
	} else if dataLength == 0 {
		return nil, fmt.Errorf("empty data at offset %d", offset)
	}

	return recordBuffer[recordStatusSize : recordStatusSize+dataLength], nil
}

func (h *FileHandler[I]) UpdateRecord(offset int64, data []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(data) > recordSize-recordStatusSize {
		return fmt.Errorf("data size is larger than max record size (%d bytes)", recordSize-recordStatusSize)
	}

	recordBuffer := make([]byte, recordSize)
	recordBuffer[0] = StatusActive
	copy(recordBuffer[recordStatusSize:], data)

	if _, err := h.dataFile.WriteAt(recordBuffer, offset); err != nil {
		return fmt.Errorf("error updating record in data file: %w", err)
	}

	return nil
}

func (h *FileHandler[I]) DeleteRecord(offset int64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, err := h.dataFile.WriteAt([]byte{StatusDeleted}, offset); err != nil {
		return fmt.Errorf("error marking record as deleted: %w", err)
	}
	return nil
}

func (h *FileHandler[I]) WriteIndexRecord(id uuid.UUID, offset int64, indexData []byte) (int64, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(indexData) > indexRecordSize-26 {
		return -1, fmt.Errorf("index data size exceeds maximum allowed size")
	}

	position, err := h.indexFile.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("error seeking to end of index file: %w", err)
	}

	record := make([]byte, indexRecordSize)
	copy(record[0:16], id[:])
	binary.LittleEndian.PutUint64(record[16:24], uint64(offset))
	binary.LittleEndian.PutUint16(record[24:26], uint16(len(indexData)))
	copy(record[26:], indexData)

	if _, err := h.indexFile.Write(record); err != nil {
		return -1, fmt.Errorf("error writing index record: %w", err)
	}

	return position, nil
}

func (h *FileHandler[I]) UpdateIndexRecord(indexOffset int64, id uuid.UUID, offset int64, indexData []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(indexData) > indexRecordSize-26 {
		return fmt.Errorf("index data size exceeds maximum allowed size")
	}

	// رفتن به موقعیت رکورد ایندکس
	if _, err := h.indexFile.Seek(indexOffset, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to index offset %d: %w", indexOffset, err)
	}

	record := make([]byte, indexRecordSize)
	copy(record[0:16], id[:])
	binary.LittleEndian.PutUint64(record[16:24], uint64(offset))
	binary.LittleEndian.PutUint16(record[24:26], uint16(len(indexData)))
	copy(record[26:], indexData)

	if _, err := h.indexFile.Write(record); err != nil {
		return fmt.Errorf("error updating index record: %w", err)
	}

	return nil
}

func (h *FileHandler[I]) ReadAllIndexRecords() (map[uuid.UUID]IndexEntry[I], error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[uuid.UUID]IndexEntry[I])

	// رفتن به ابتدای فایل ایندکس
	if _, err := h.indexFile.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("error seeking to start of index file: %w", err)
	}

	record := make([]byte, indexRecordSize)
	currentOffset := int64(0)

	for {
		n, err := h.indexFile.Read(record)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading index record: %w", err)
		}

		if n < indexRecordSize {
			break
		}

		// نادیده گرفتن رکوردهای خالی
		if record[0] == 0 {
			currentOffset += indexRecordSize
			continue
		}

		// خواندن UUID
		id, err := uuid.FromBytes(record[0:16])
		if err != nil {
			log.Printf("Error parsing UUID at offset %d: %v", currentOffset, err)
			currentOffset += indexRecordSize
			continue
		}

		// خواندن آفست داده
		dataOffset := int64(binary.LittleEndian.Uint64(record[16:24]))

		// خواندن اندازه داده
		dataSize := binary.LittleEndian.Uint16(record[24:26])

		// خواندن داده
		data := make([]byte, dataSize)
		copy(data, record[26:26+dataSize])

		// Unmarshal داده ایندکس
		var indexData I
		if err := json.Unmarshal(data, &indexData); err != nil {
			log.Printf("Error unmarshaling index data for ID %s at offset %d: %v", id, currentOffset, err)
			currentOffset += indexRecordSize
			continue
		}

		// اضافه کردن به نتیجه
		result[id] = IndexEntry[I]{
			Offset:      dataOffset,
			IndexOffset: currentOffset,
			IndexData:   indexData,
		}

		currentOffset += indexRecordSize
	}

	return result, nil
}

type Manager[T collectionItem, I collectionItem] struct {
	fh           *FileHandler[I]
	mu           sync.RWMutex
	primaryIndex map[uuid.UUID]IndexEntry[I]
	closed       bool
}

func New[T collectionItem, I collectionItem]() (*Manager[T, I], error) {
	fh, err := NewFileHandler[I]()
	if err != nil {
		return nil, fmt.Errorf("failed to create file handler: %w", err)
	}

	manager := &Manager[T, I]{
		fh:           fh,
		primaryIndex: make(map[uuid.UUID]IndexEntry[I]),
	}

	if err := manager.loadPrimaryIndex(); err != nil {
		return nil, fmt.Errorf("failed to load primary index: %w", err)
	}

	return manager, nil
}

func (m *Manager[T, I]) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}
	m.closed = true

	return m.fh.Close()
}

// loadPrimaryIndex ایندکس اصلی را از دیسک بارگذاری می‌کند
func (m *Manager[T, I]) loadPrimaryIndex() error {
	// بررسی وجود فایل ایندکس
	if _, err := os.Stat(filepath.Join(dirName, "index.db")); os.IsNotExist(err) {
		return m.rebuildIndex()
	}

	// خواندن تمام رکوردهای ایندکس
	indexMap, err := m.fh.ReadAllIndexRecords()
	if err != nil {
		return fmt.Errorf("error reading index records: %w", err)
	}

	if len(indexMap) == 0 {
		return m.rebuildIndex()
	}

	// اختصاص مستقیم ایندکس خوانده شده
	m.primaryIndex = indexMap

	log.Printf("Loaded %d entries from primary index", len(m.primaryIndex))
	return nil
}

func (m *Manager[T, I]) rebuildIndex() error {
	m.primaryIndex = make(map[uuid.UUID]IndexEntry[I])
	fileInfo, err := m.fh.dataFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting data file info: %w", err)
	}
	fileSize := fileInfo.Size()

	if fileSize == 0 {
		log.Println("Data file is empty, no index to rebuild")
		return nil
	}

	// پاک کردن فایل ایندکس قبل از بازسازی
	if err := m.fh.indexFile.Truncate(0); err != nil {
		return fmt.Errorf("error truncating index file: %w", err)
	}
	if _, err := m.fh.indexFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to start of index file: %w", err)
	}

	for offset := int64(0); offset < fileSize; offset += recordSize {
		recordBuffer := make([]byte, recordSize)
		n, err := m.fh.dataFile.ReadAt(recordBuffer, offset)
		if err != nil && err != io.EOF {
			log.Printf("Error reading record at offset %d: %v", offset, err)
			continue
		}

		if n == 0 {
			continue
		}

		if recordBuffer[0] == StatusDeleted {
			continue
		}

		dataLength := bytes.IndexByte(recordBuffer[recordStatusSize:], 0)
		if dataLength == -1 {
			dataLength = recordSize - recordStatusSize
		}

		if dataLength == 0 {
			continue
		}

		data := recordBuffer[recordStatusSize : recordStatusSize+dataLength]

		var dataItem T
		if err := json.Unmarshal(data, &dataItem); err != nil {
			log.Printf("Error unmarshaling data at offset %d: %v", offset, err)
			continue
		}

		indexItem, err := createIndexItem[T, I](dataItem)
		if err != nil {
			log.Printf("Error creating index item from data at offset %d: %v", offset, err)
			continue
		}

		id := dataItem.GetID()
		if id != uuid.Nil {
			indexData, err := json.Marshal(indexItem)
			if err != nil {
				log.Printf("Error marshaling index item for ID %s: %v", id, err)
				continue
			}

			// نوشتن رکورد ایندکس و دریافت موقعیت آن
			indexOffset, err := m.fh.WriteIndexRecord(id, offset, indexData)
			if err != nil {
				log.Printf("Error writing index record for ID %s: %v", id, err)
				continue
			}

			m.primaryIndex[id] = IndexEntry[I]{
				Offset:      offset,
				IndexOffset: indexOffset,
				IndexData:   indexItem,
			}
		}
	}

	log.Printf("Rebuilt primary index with %d entries", len(m.primaryIndex))
	return nil
}

func createIndexItem[T, I any](dataItem T) (I, error) {
	var zero I
	dataValue := reflect.ValueOf(dataItem)
	indexType := reflect.TypeOf(zero)

	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	if dataValue.Kind() != reflect.Struct {
		return zero, fmt.Errorf("data type must be a struct or pointer to struct")
	}

	var isIndexPointer bool
	var indexElemType reflect.Type

	if indexType.Kind() == reflect.Ptr {
		isIndexPointer = true
		indexElemType = indexType.Elem()
	} else {
		indexElemType = indexType
	}

	if indexElemType.Kind() != reflect.Struct {
		return zero, fmt.Errorf("index type must be a struct or pointer to struct")
	}

	var indexValue reflect.Value
	if isIndexPointer {
		indexValue = reflect.New(indexElemType)
	} else {
		indexValue = reflect.New(indexElemType).Elem()
	}

	indexFields := make(map[string]reflect.StructField)
	for i := 0; i < indexElemType.NumField(); i++ {
		field := indexElemType.Field(i)
		indexFields[field.Name] = field
	}

	dataType := dataValue.Type()
	for i := 0; i < dataType.NumField(); i++ {
		dataField := dataType.Field(i)
		indexTag := dataField.Tag.Get("index")
		if indexTag != "true" {
			continue
		}

		if indexField, ok := indexFields[dataField.Name]; ok {
			dataFieldValue := dataValue.Field(i)
			var indexFieldValue reflect.Value

			if isIndexPointer {
				indexFieldValue = indexValue.Elem().FieldByName(dataField.Name)
			} else {
				indexFieldValue = indexValue.FieldByName(dataField.Name)
			}

			if dataFieldValue.IsValid() && indexFieldValue.IsValid() &&
				dataFieldValue.Type().AssignableTo(indexField.Type) {
				indexFieldValue.Set(dataFieldValue)
			} else {
				log.Printf("Warning: Cannot assign field %s from type %v to type %v",
					dataField.Name, dataFieldValue.Type(), indexField.Type)
			}
		}
	}

	if isIndexPointer {
		return indexValue.Interface().(I), nil
	}
	return indexValue.Interface().(I), nil
}

func (m *Manager[T, I]) Create(item T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero T
	if m.closed {
		return zero, fmt.Errorf("manager is closed")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return zero, fmt.Errorf("error generating UUID v7: %w", err)
	}
	item.SetID(id)

	data, err := json.Marshal(item)
	if err != nil {
		return zero, fmt.Errorf("error marshaling item: %w", err)
	}

	offset, err := m.fh.WriteRecord(data)
	if err != nil {
		return zero, fmt.Errorf("error writing record: %w", err)
	}

	indexItem, err := createIndexItem[T, I](item)
	if err != nil {
		m.fh.DeleteRecord(offset)
		return zero, fmt.Errorf("failed to create index item: %w", err)
	}

	indexData, err := json.Marshal(indexItem)
	if err != nil {
		m.fh.DeleteRecord(offset)
		return zero, fmt.Errorf("failed to marshal index item: %w", err)
	}

	// نوشتن رکورد ایندکس و دریافت موقعیت آن
	indexOffset, err := m.fh.WriteIndexRecord(id, offset, indexData)
	if err != nil {
		m.fh.DeleteRecord(offset)
		return zero, fmt.Errorf("failed to write index record: %w", err)
	}

	m.primaryIndex[id] = IndexEntry[I]{
		Offset:      offset,
		IndexOffset: indexOffset,
		IndexData:   indexItem,
	}

	return item, nil
}

func (m *Manager[T, I]) Update(item T) (T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var zero T
	id := item.GetID()
	entry, ok := m.primaryIndex[id]
	if !ok {
		return zero, fmt.Errorf("item with ID %s does not exist", id.String())
	}

	indexItem, err := createIndexItem[T, I](item)
	if err != nil {
		return zero, fmt.Errorf("failed to create index item: %w", err)
	}

	data, err := json.Marshal(item)
	if err != nil {
		return zero, fmt.Errorf("error marshaling item: %w", err)
	}

	if err := m.fh.UpdateRecord(entry.Offset, data); err != nil {
		return zero, fmt.Errorf("error updating record: %w", err)
	}

	indexData, err := json.Marshal(indexItem)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal index item: %w", err)
	}

	// به‌روزرسانی رکورد ایندکس موجود به جای نوشتن رکورد جدید
	if err := m.fh.UpdateIndexRecord(entry.IndexOffset, id, entry.Offset, indexData); err != nil {
		return zero, fmt.Errorf("failed to update index record: %w", err)
	}

	m.primaryIndex[id] = IndexEntry[I]{
		Offset:      entry.Offset,
		IndexOffset: entry.IndexOffset,
		IndexData:   indexItem,
	}

	return item, nil
}

func (m *Manager[T, I]) Delete(id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.primaryIndex[id]
	if !ok {
		return fmt.Errorf("item with ID %s not found", id)
	}

	if err := m.fh.DeleteRecord(entry.Offset); err != nil {
		return err
	}

	// حذف رکورد ایندکس با نوشتن یک رکورد خالی
	emptyRecord := make([]byte, indexRecordSize)
	if _, err := m.fh.indexFile.WriteAt(emptyRecord, entry.IndexOffset); err != nil {
		return fmt.Errorf("error deleting index record: %w", err)
	}

	delete(m.primaryIndex, id)

	return nil
}

func (m *Manager[T, I]) Read(id uuid.UUID) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var zero T
	entry, ok := m.primaryIndex[id]
	if !ok {
		return zero, fmt.Errorf("item not found with ID: %s", id)
	}

	data, err := m.fh.ReadRecord(entry.Offset)
	if err != nil {
		return zero, fmt.Errorf("error reading item from disk: %w", err)
	}

	var loadedItem T
	if err := json.Unmarshal(data, &loadedItem); err != nil {
		return zero, fmt.Errorf("error unmarshaling item: %w", err)
	}

	loadedItem.SetID(id)
	return loadedItem, nil
}

func (m *Manager[T, I]) ReadIndex(id uuid.UUID) (I, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var zero I
	entry, exists := m.primaryIndex[id]
	if !exists {
		return zero, fmt.Errorf("index not found with ID: %s", id)
	}
	return entry.IndexData, nil
}

func (m *Manager[T, I]) IsClosed() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.closed
}

func (m *Manager[T, I]) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.primaryIndex)
}

func (m *Manager[T, I]) DebugInfo() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return map[string]interface{}{
		"primary_index_size": len(m.primaryIndex),
		"is_closed":          m.closed,
	}
}

func (m *Manager[T, I]) PrintDebugInfo() {
	info := m.DebugInfo()
	log.Printf("Debug Info - Primary Index: %d, Closed: %v",
		info["primary_index_size"],
		info["is_closed"])
}

func (m *Manager[T, I]) GetIndexEntry(id uuid.UUID) (IndexEntry[I], error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.primaryIndex[id]
	var zero IndexEntry[I]
	if !exists {
		return zero, fmt.Errorf("index entry not found with ID: %s", id)
	}

	return entry, nil
}

func (m *Manager[T, I]) GetAllIndexEntries() map[uuid.UUID]IndexEntry[I] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[uuid.UUID]IndexEntry[I], len(m.primaryIndex))
	for k, v := range m.primaryIndex {
		result[k] = v
	}

	return result
}
