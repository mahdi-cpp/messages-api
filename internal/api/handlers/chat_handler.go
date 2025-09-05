package handlers

import (
	"fmt"
	"net/http"

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

	fmt.Println(c.Param("chatId"))

	chatID := c.Param("chatId")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	fmt.Println("chatID:", chatID)
	fmt.Println("userID:", userID)

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
