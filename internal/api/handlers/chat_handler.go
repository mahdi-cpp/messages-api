package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/application"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

type ChatHandler struct {
	appManager *application.Manager
}

func NewChatHandler(appManager *application.Manager) *ChatHandler {
	return &ChatHandler{
		appManager: appManager,
	}
}

// Create godoc
// @Router /api/chats/{chatId}/users/{userId}/createChat [put]
func (h *ChatHandler) Create(c *gin.Context) {

	fmt.Println("chat create handler")

	var request *chat.Chat
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.AbortWithRequestInvalid(c)
		return
	}

	err := h.appManager.CreateChat(request)
	if err != nil {
		return
	}
}

func (h *ChatHandler) Read(c *gin.Context) {

}

func (h *ChatHandler) Update(c *gin.Context) {

}

func (h *ChatHandler) Delete(c *gin.Context) {

}

func (h *ChatHandler) List(c *gin.Context) {}
