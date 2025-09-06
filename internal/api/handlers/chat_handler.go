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

	chatID := c.Param("id")

	chat1, err := h.appManager.ReadChat(chatID)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": chat1})
}

func (h *ChatHandler) ReadAll(c *gin.Context) {

	var request chat.SearchOptions
	if err := c.ShouldBindQuery(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("title:", request.Title)
	fmt.Println("offset:", request.Offset)
	fmt.Println("limit:", request.Limit)
	if request.IsVerified != nil {
		fmt.Println("verified:", *request.IsVerified)
	}

	chats, err := h.appManager.ReadAllChats(&request)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

func (h *ChatHandler) Update(c *gin.Context) {

	chatID := c.Param("id")

	var request chat.UpdateOptions
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body: " + err.Error()})
		return
	}

	//for _, member := range request.Members {
	//	fmt.Println("members:", member.UserID)
	//}

	err := h.appManager.UpdateChat(chatID, request)
	if err != nil {
		c.JSON(http.StatusNotModified, gin.H{"message": "failed update chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully updated"})
}

func (h *ChatHandler) BuckUpdate(c *gin.Context) {
	var request BuckUpdateChats

	// Use c.ShouldBindJSON to bind the request body.
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body: " + err.Error()})
		return
	}

	for _, id := range request.Ids {
		fmt.Println(id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chats updated successfully"})
}

func (h *ChatHandler) Delete(c *gin.Context) {
	// Get the ID from the URL path
	id := c.Param("id")

	// The 'id' variable will contain "12" from the URL
	fmt.Printf("Deleting chat with ID: %s\n", id)

	// ... your chat deletion logic here

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Chat with ID %s deleted", id)})
}

type BuckUpdateChats struct {
	Ids []string `json:"ids"`
}

type BuckDeleteChats struct {
	Ids []string `json:"ids"`
}

func (h *ChatHandler) BuckDelete(c *gin.Context) {
	var request BuckDeleteChats

	// Use c.ShouldBindJSON to bind the request body.
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body: " + err.Error()})
		return
	}

	for _, id := range request.Ids {
		fmt.Println(id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chats deleted successfully"})
}

func (h *ChatHandler) List(c *gin.Context) {}
