package collection_manager_db

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
)

func TestCreateMessages(t *testing.T) {

	db, err := New[*message.Message]()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	all, err := db.ReadAll()
	if err != nil {
		return
	}

	for _, msg := range all {
		fmt.Println(msg.Caption)
	}

	var count int

	for i := 0; i < 1000000; i++ {
		msg := &message.Message{
			ID:      uuid.New(),
			Caption: "caption of message number  " + strconv.Itoa(i),
			Medias: []*message.Media{
				{
					ID:          uuid.New(),
					MimeType:    "image/jpeg",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          uuid.New(),
					MimeType:    "image/jpeg",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          uuid.New(),
					MimeType:    "video/mp4",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          uuid.New(),
					MimeType:    "video/mp4",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := db.Create(msg)
		if err != nil {
			t.Fatal(err)
		}
		count++
	}

	fmt.Printf("Created %d messages.\n", count)
}

func TestReadAllMessages(t *testing.T) {

	db, err := New[*message.Message]()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	all, err := db.ReadAll()
	if err != nil {
		return
	}

	fmt.Printf("Read %d messages.\n", len(all))

	//for _, msg := range all {
	//	fmt.Println(msg.Caption)
	//}
}

func TestReadMessages(t *testing.T) {

	db, err := New[*message.Message]()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	start := time.Now()

	id, err := uuid.Parse("01995280-a910-742c-aaac-4ccb911e216a")
	if err != nil {
		t.Fatal(err)
	}

	msg, err := db.Read(id)
	if err != nil {
		t.Fatal(err)
	}

	duration := time.Since(start)
	fmt.Println(duration)

	fmt.Println(msg.Caption)
}
