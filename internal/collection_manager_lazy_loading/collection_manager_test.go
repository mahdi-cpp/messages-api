package collection_manager_lazy_loading

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
)

// A place to store the created UUIDs for the read test
//var createdUUIDs []uuid.UUID

func TestCreateMessages(t *testing.T) {

	db, err := New[*message.Message]()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Use a new slice for each test run to avoid conflicts
	//createdUUIDs = make([]uuid.UUID, 0, 60000)

	v7, err := uuid.NewV7()
	if err != nil {
		return
	}

	var count int
	for i := 0; i < 100; i++ {
		msg := &message.Message{
			// DO NOT set the ID here. Let the Create function handle it.
			Caption:  "caption  " + strconv.Itoa(i),
			IsEdited: true,
			IsPinned: true,
			Medias: []*message.Media{
				{
					ID:          v7,
					MimeType:    "image/jpeg",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "video/mp4",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "video/webm",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "image/jpeg",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "video/mp4",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "video/webm",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				}, {
					ID:          v7,
					MimeType:    "video/mp4",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "video/webm",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "image/jpeg",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
				{
					ID:          v7,
					MimeType:    "video/mp4",
					Width:       1000,
					Height:      1000,
					Orientation: "portrait",
				},
			},
			Voice: &message.Voice{
				ID:       v7,
				MimeType: "audio/opus",
				Duration: 5000,
				FileSize: 25000,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := db.Create(msg)
		if err != nil {
			t.Fatal(err)
		}
		// Store the new UUID returned by the Create function
		//createdUUIDs = append(createdUUIDs, createdMsg.GetID())
		count++
	}

	fmt.Printf("Created %d messages.\n", count)
}

func TestReadMessages(t *testing.T) {

	//jsonCollection, err := collection_manager.New[*message.Message]("/app/tmp/messages/metadata")
	//if err != nil {
	//	return
	//}

	db, err := New[*message.Message]()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	all := db.fh.GetAllUUIDs()
	fmt.Printf("Read %d messages.\n", len(all))
	fmt.Printf("Cpation: %s\n", all[0].String())
	if len(all) == 0 {
		t.Fatal("No messages found")
	}

	start := time.Now()

	for i := 0; i < 100; i++ {
		id := all[i]
		//fmt.Printf("id: %s \n", id.String())

		msg, err := db.Read(id)
		if err != nil {
			t.Fatalf("Error reading message with UUID %s: %v", id, err)
		}
		fmt.Println(msg.Caption)

		//_, err = jsonCollection.Create(msg)
		//if err != nil {
		//	t.Fatal(err)
		//}
	}

	duration := time.Since(start)
	fmt.Println("single  read duration: ", duration)

}
