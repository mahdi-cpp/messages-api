package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/application"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/collections/message"
	"github.com/mahdi-cpp/messages-api/internal/helpers"
)

type MessageHandler struct {
	appManager *application.Manager
}

func NewMessageHandler(appManager *application.Manager) *MessageHandler {
	return &MessageHandler{
		appManager: appManager,
	}
}

// Create godoc
// @Router /api/messages post
func (h *MessageHandler) Create(c *gin.Context) {

	var request *message.Message
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.AbortWithRequestInvalid(c)
		return
	}

	fmt.Println(request.ChatID)

	newMessage, err := h.appManager.MessageCreate(request)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, newMessage)
}

// ShowAccount godoc
//
//	@Summary		Show an account
//	@Description	get string by ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	model.Account
//	@Failure		400	{object}	http util.HTTPError
//	@Failure		404	{object}	http util.HTTPError
//	@Failure		500	{object}	http util.HTTPError
//	@Router			/chats/{id} [get]
func (h *MessageHandler) Read(c *gin.Context) {

	chatID := c.Param("id")

	chat1, err := h.appManager.ReadChat(chatID)
	if err != nil {
		return
	}

	fmt.Println(chat1.Title)

	c.JSON(http.StatusOK, chat1)
}

func (h *MessageHandler) ReadAll(c *gin.Context) {

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

	fmt.Println("ReadAll match:", len(chats))

	c.JSON(http.StatusOK, chats)
}

func (h *MessageHandler) Update(c *gin.Context) {

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

func (h *MessageHandler) BuckUpdate(c *gin.Context) {
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

func (h *MessageHandler) Delete(c *gin.Context) {
	// Get the ID from the URL path
	id := c.Param("id")

	// The 'id' variable will contain "12" from the URL
	fmt.Printf("Deleting chat with ID: %s\n", id)

	// ... your chat deletion logic here

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Chat with ID %s deleted", id)})
}

func (h *MessageHandler) BuckDelete(c *gin.Context) {
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
