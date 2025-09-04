package yourService

import (
	"context"
	"fmt"

	telegram2 "github.com/mahdi-cpp/messages-api/internal/telegram/models"
)

func UpdateChatAbout(context context.Context, t *telegram2.EditChatAboutRequest) (interface{}, error) {
	return t.ChatID, fmt.Errorf("")
}

func DeleteChat(ctx context.Context, id string) error {
	return fmt.Errorf("")
}

func CreateChat(ctx context.Context, t *telegram2.CreateChatRequest) (interface{}, error) {
	return nil, fmt.Errorf("")
}

func GetChats(ctx context.Context, ds []string) (interface{}, error) {
	return nil, fmt.Errorf("")
}

func UpdateChatAdminRights(ctx context.Context, id string, id2 string, t *telegram2.ChatAdminRights) (interface{}, error) {
	return nil, fmt.Errorf("")
}

func UpdateChatBannedRights(ctx context.Context, id string, id2 string, t *telegram2.ChatBannedRights) (interface{}, error) {
	return nil, fmt.Errorf("")
}
