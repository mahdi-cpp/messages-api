package message

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collection_manager_gemini_v2"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

func TestReadMessages(t *testing.T) {

	var err error
	var messagesDirectory = config.GetPath("messages_test1")
	collectionManager, err := collection_manager_gemini_v2.New[*Message](messagesDirectory)
	if err != nil {
		t.Error(err)
	}

	all, err := collectionManager.ReadAll()
	if err != nil {
		return
	}

	for _, msg := range all {
		if msg.Medias != nil {
			fmt.Println(msg.Content, msg.Medias[0].MimeType)
		}
	}
}

func TestMessageCreate(t *testing.T) {

	var err error
	var messagesDirectory = config.GetPath("messages_test1")
	collectionManager, err := collection_manager_gemini_v2.New[*Message](messagesDirectory)
	if err != nil {
		t.Error(err)
	}

	all, err := collectionManager.ReadAll()
	if err != nil {
		return
	}

	for _, msg := range all {
		fmt.Println(msg.Music)
	}

	msg := &Message{

		ID:      uuid.New(),
		ChatID:  uuid.New(),
		UserID:  uuid.New(),
		Content: "test content",

		//Medias: []*Media{
		//	{
		//		ID:       uuid.New(),
		//		Width:    100,
		//		Height:   200,
		//		MimeType: "image/jpeg",
		//	},
		//	{
		//		ID:       uuid.New(),
		//		Width:    100,
		//		Height:   200,
		//		MimeType: "video/mp4",
		//	},
		//},

		Music: &Music{
			ID:       uuid.New(),
			MimeType: "audio/mpeg",
			Duration: 12,
			FileSize: 80000,
		},
		Version: "1",
	}

	create, err := collectionManager.Create(msg)
	if err != nil {
		t.Errorf("Error creating message: %s", err)
	}

	fmt.Println(create)

}
