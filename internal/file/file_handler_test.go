package collection_manager_db

import (
	"fmt"
	"log"
	"testing"
)

//func TestNewFileHandler(t *testing.T) {
//
//	start := time.Now()
//
//	fileHandler, err := NewFileHandler()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer fileHandler.Close()
//
//	var count = 60000
//
//	//for i := 0; i < count; i++ {
//	//	err = fileHandler.WriteRecord(i, []byte("row 0123456789"+strconv.Itoa(i)))
//	//	if err != nil {
//	//		log.Fatal(err)
//	//	}
//	//}
//
//	for i := 0; i < count; i++ {
//		_, err := fileHandler.readRecord(i)
//		if err != nil {
//			log.Fatal(err)
//		}
//		//fmt.Printf("row : %s\n", data)
//	}
//
//	duration := time.Since(start)
//	fmt.Println(duration)
//}

func TestNewFileHandler(t *testing.T) {

	handler, err := NewFileHandler()
	if err != nil {
		log.Fatal(err)
	}
	defer handler.Close()

	// مثال: نوشتن یک رکورد جدید
	data := []byte("First data save")
	newUUID, err := handler.WriteRecord(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("داده با موفقیت ذخیره شد. UUID: %s\n", newUUID)

	// مثال: خواندن همان رکورد با UUID
	readData, err := handler.ReadRecord(newUUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read datas: %s\n", readData)
}
