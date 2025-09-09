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

// Create
// @Summary     create a new chat
// @Description Creates a new chat instance. The request body should contain a list of members to include in the chat.
// @Tags        chat
// @Accept      json
// @Produce     json
// @Param       request body chat.Chat true "Chat creation request body"
// @Success     201 {object} chat.Chat "Chat created successfully"
// @Failure     400 {string
// @Router      /chats [post]
func (h *ChatHandler) Create(c *gin.Context) {

	var request *chat.Chat
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.AbortWithRequestInvalid(c)
		return
	}

	fmt.Println(request.Members)

	newChat, err := h.appManager.ChatCreate(request)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, newChat)
}

// ReadChat
// @Description Retrieves a single chat instance by its unique ID.
// @Tags chat
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} chat.Chat "Chat retrieved successfully"
// @Failure 404 {string} string "Chat not found"
// @Failure 500 {string} string "Internal server error"
// @Router /chats/{id} [get]
func (h *ChatHandler) ReadChat(c *gin.Context) {

	chatID := c.Param("id")

	chat1, err := h.appManager.ReadChat(chatID)
	if err != nil {
		return
	}

	fmt.Println(chat1.Title)

	c.JSON(http.StatusOK, chat1)
}

// Read
// @Summary Get a list of chats
// @Description Retrieves a list of chats, with optional search, pagination, and filtering.
// @Tags chat
// @Accept json
// @Produce json
// @Param title query string false "Filter by chat title"
// @Param offset query int false "Pagination offset"
// @Param limit query int false "Pagination limit"
// @Param isVerified query boolean false "Filter by verification status"
// @Success 200 {array} chat.Chat "List of chats retrieved successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /chats [get]
func (h *ChatHandler) Read(c *gin.Context) {

	var request chat.SearchOptions
	if err := c.ShouldBindQuery(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if request.ChatID == "" { //read all chats with Chat SearchOptions
		h.readAllChats(c, &request)
	} else if request.ChatID != "" {
		h.readSingleChat(c, request.ChatID)
	}
}

// Private helper methods
func (h *ChatHandler) readAllChats(c *gin.Context, options *chat.SearchOptions) {

	chats, err := h.appManager.ReadAllChats(options)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("title:", options.Title)
	fmt.Println("page:", options.Page)
	fmt.Println("size:", options.Size)
	if options.IsVerified != nil {
		fmt.Println("verified:", *options.IsVerified)
	}

	fmt.Println("ReadAllMessages match:", len(chats))
	c.JSON(http.StatusOK, chats)
}

func (h *ChatHandler) readSingleChat(c *gin.Context, chatID string) {

	fmt.Println("readSingleChat", chatID)
	readChat, err := h.appManager.ReadChat(chatID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readChat)
}

//func (h *ChatHandler) ReadChatMessage(c *gin.Context) {
//	chatID := c.Param("chatId")
//	messageID := c.Param("messageId")
//
//	chatManager, err := h.appManager.GetChatManager(chatID)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
//		return
//	}
//
//	message, err := chatManager.ReadMessage(messageID)
//	if err != nil {
//		return
//	}
//	c.JSON(http.StatusOK, message)
//}
//
//func (h *ChatHandler) ReadChatMessages(c *gin.Context) {
//
//}

// Update
// @Summary update an existing chat
// @Description Updates an existing chat's properties, such as its members or title.
// @Tags chat
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body chat.UpdateOptions true "Chat update options"
// @Success 200 {object} object "Chat updated successfully"
// @Failure 304 {string} string "Not Modified"
// @Failure 400 {string} string "Invalid JSON body"
// @Router /chats/{id} [put]
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

// BuckUpdate
// @Summary Bulk update multiple chats
// @Description Updates multiple chats based on a list of IDs provided in the request body.
// @Tags chat
// @Accept json
// @Produce json
// @Param request body BuckUpdateChats true "List of chat IDs to update"
// @Success 200 {object} object "Chats updated successfully"
// @Failure 400 {string} string "Invalid JSON body"
// @Router /chats/bulk-update [put]
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

// Delete
// @Summary delete a single chat by ID
// @Description Deletes a chat instance by its unique ID.
// @Tags chat
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} object "Chat deleted successfully"
// @Router /chats/{id} [delete]
func (h *ChatHandler) Delete(c *gin.Context) {
	// Get the ID from the URL path
	id := c.Param("id")

	// The 'id' variable will contain "12" from the URL
	fmt.Printf("Deleting chat with ID: %s\n", id)

	// ... your chat deletion logic here

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Chat with ID %s deleted", id)})
}

// BuckDelete
// @Summary Bulk delete multiple chats
// @Description Deletes multiple chats based on a list of IDs provided in the request body.
// @Tags chat
// @Accept json
// @Produce json
// @Param request body BuckDeleteChats true "List of chat IDs to delete"
// @Success 200 {object} object "Chats deleted successfully"
// @Failure 400 {string} string "Invalid JSON body"
// @Router /chats/bulk-delete [delete]
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
