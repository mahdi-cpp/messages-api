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
// @Router /api/chats post
func (h *ChatHandler) Create(c *gin.Context) {

	var request *chat.Chat
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.AbortWithRequestInvalid(c)
		return
	}

	fmt.Println(request.Members)

	err := h.appManager.CreateChat(request)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat": "created"})
}

func (h *ChatHandler) Read(c *gin.Context) {

}

func (h *ChatHandler) ReadAll(c *gin.Context) {

	var request chat.SearchOptions
	if err := c.ShouldBindQuery(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("title:", request.Title)

	if request.IsVerified != nil {
		fmt.Println("verified:", *request.IsVerified)
	}

	fmt.Println("offset:", request.Offset)
	fmt.Println("limit:", request.Limit)

	c.JSON(http.StatusOK, gin.H{"message": "Search executed successfully"})
}

func (h *ChatHandler) Update(c *gin.Context) {

}

func (h *ChatHandler) Delete(c *gin.Context) {

}

func (h *ChatHandler) List(c *gin.Context) {}
