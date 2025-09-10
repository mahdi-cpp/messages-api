package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/application"
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

	fmt.Println("MessageHandler create")

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
//func (h *MessageHandler) Read(c *gin.Context) {
//
//	var request message.SearchOptions
//	if err := c.ShouldBindQuery(&request); err != nil {
//		fmt.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	if request.MessageID == "" { //read all message with SearchOptions
//		selectedMessages, err := h.appManager.ReadAllMessages(&request)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//
//		c.JSON(http.StatusOK, selectedMessages)
//
//	} else { //read single message
//
//		c.JSON(http.StatusNotFound, gin.H{"messageId not found": request.MessageID})
//	}
//}
//
//func (h *MessageHandler) ReadAll(c *gin.Context) {
//
//	fmt.Println("MessageHandler readAll")
//
//	var with message.SearchOptions
//	if err := c.ShouldBindQuery(&with); err != nil {
//		fmt.Println(err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	fmt.Println("with:", with)
//
//	fmt.Println("offset:", with.Page)
//	fmt.Println("limit:", with.Size)
//
//	messages, err := h.appManager.ReadAllMessages(&with)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	fmt.Println("ReadAllMessages match:", len(messages))
//
//	c.JSON(http.StatusOK, messages)
//}

func (h *MessageHandler) Read(c *gin.Context) {

	var request message.SearchOptions
	if err := c.ShouldBindQuery(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if request.MessageID == uuid.Nil { //read all messages with Message SearchOptions
		h.readAllMessage(c, &request)
	} else if request.MessageID != uuid.Nil {
		h.readSingleMessage(c, request.ChatID, request.MessageID)
	}
}

func (h *MessageHandler) readAllMessage(c *gin.Context, options *message.SearchOptions) {
	fmt.Println("readAllMessage")

	selectedMessages, err := h.appManager.ReadAllMessages(options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, selectedMessages)
}

func (h *MessageHandler) readSingleMessage(c *gin.Context, chatID, messageId uuid.UUID) {

	fmt.Println("readSingleMessage", chatID)
	chatManager, err := h.appManager.GetChatManager(chatID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	readMessage, err := chatManager.ReadMessage(messageId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readMessage)
}

func (h *MessageHandler) Update(c *gin.Context) {

	var request message.UpdateOptions
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body: " + err.Error()})
		return
	}

	chatManager, err := h.appManager.GetChatManager(request.ChatID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	messageUpdated, err := chatManager.UpdateMessage(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messageUpdated)
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
	// Read the ID from the URL path
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
