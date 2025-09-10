package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/telegram"
	telegram2 "github.com/mahdi-cpp/messages-api/internal/telegram/models"
)

// GetChats godoc
// @Summary Read a list of chats by ID
// @Description Retrieves a list of chat objects based on their IDs.
// @Accept  json
// @Produce  json
// @Param chatIDs query string true "A comma-separated list of chat IDs to retrieve"
// @Success 200 {object} GetChatsResponse
// @Failure 400 {object} map[string]string "Invalid input"
// @Router /chats [get]
func GetChats(c *gin.Context) {
	// این خط برای دریافت پارامترهای URL استفاده می‌شود
	chatIDsStr := c.Query("chatIDs")
	if chatIDsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chatIDs query parameter is required"})
		return
	}

	// تبدیل رشته ID به slice
	chatIDs := strings.Split(chatIDsStr, ",")

	// فراخوانی لایه سرویس برای دریافت چت‌ها
	chats, err := yourService.GetChats(c.Request.Context(), chatIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

// CreateChat godoc
// @Summary Create a new chat
// @Description Creates a new chat with a title and initial users.
// @Accept  json
// @Produce  json
// @Param request body CreateChatRequest true "Chat creation request"
// @Success 201 {object} CreateChatResponse
// @Failure 400 {object} map[string]string "Invalid request body"
// @Router /chats [post]
func CreateChat(c *gin.Context) {
	var req telegram2.CreateChatRequest

	// ShouldBindJSON برای اتصال (bind) بدنه JSON درخواست به struct استفاده می‌شود
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// فراخوانی لایه سرویس برای ایجاد چت
	newChat, err := yourService.CreateChat(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, newChat)
}

// DeleteChat godoc
// @Summary Delete a chat
// @Description Deletes a chat by its ID.
// @Accept  json
// @Produce  json
// @Param id path string true "Chat ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string "Chat not found"
// @Router /chats/{id} [delete]
func DeleteChat(c *gin.Context) {

	// دریافت پارامتر path از URL
	chatID := c.Param("id")

	// فراخوانی لایه سرویس برای حذف چت
	if err := yourService.DeleteChat(c.Request.Context(), chatID); err != nil {
		// بررسی نوع خطا
		if err.Error() == "not found" { // این بخش بستگی به نوع خطای برگردانده شده از سرویس دارد
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		}
		return
	}

	// ارسال پاسخ موفقیت‌آمیز بدون محتوا
	c.Status(http.StatusNoContent)
}

// EditChatAbout godoc
// @Summary Edit a chat's "about" section
// @Description Updates the "about" section for a specific chat.
// @Accept  json
// @Produce  json
// @Param id path string true "Chat ID"
// @Param request body EditChatAboutRequest true "Chat about update request"
// @Success 200 {object} Chat "Updated chat object"
// @Failure 400 {object} map[string]string "Invalid request body or parameters"
// @Failure 404 {object} map[string]string "Chat not found"
// @Router /chats/{id}/about [put]
func EditChatAbout(c *gin.Context) {

	// 1. دریافت ID چت از URL
	chatID := c.Param("id")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})
		return
	}

	// 2. اتصال بدنه JSON درخواست به struct
	var req telegram2.EditChatAboutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. (اختیاری) اعتبارسنجی اضافی
	// مطمئن شوید که chatId در بدنه و URL یکسان است یا از یکی استفاده کنید.
	// در اینجا، از ChatID در URL به عنوان منبع اصلی استفاده می‌کنیم.
	if req.ChatID != "" && req.ChatID != chatID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mismatched chat ID in URL and body"})
		return
	}
	req.ChatID = chatID // اطمینان از یکسان بودن ID

	// 4. فراخوانی لایه سرویس برای به‌روزرسانی
	updatedChat, err := yourService.UpdateChatAbout(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat about"})
		}
		return
	}

	// 5. ارسال پاسخ موفقیت‌آمیز
	c.JSON(http.StatusOK, updatedChat)
}

// UpdateChatAdminRights godoc
// @Summary Update admin rights for a user in a chat
// @Description Updates the administrative rights of a user within a chat.
// @Accept  json
// @Produce  json
// @Param chatId path string true "ID of the chat"
// @Param userId path string true "ID of the user whose rights are being updated"
// @Param request body ChatAdminRights true "The admin rights to apply"
// @Success 200 {object} Chat "The updated chat object"
// @Failure 400 {object} map[string]string "Invalid request or parameters"
// @Failure 404 {object} map[string]string "Chat or user not found"
// @Failure 403 {object} map[string]string "Permission denied"
// @Router /chats/{chatId}/admins/{userId}/rights [put]
func UpdateChatAdminRights(c *gin.Context) {

	// دریافت شناسه های چت و کاربر از پارامترهای URL
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

	// اتصال بدنه JSON درخواست به ساختار ChatAdminRights
	var adminRights telegram2.ChatAdminRights
	if err := c.ShouldBindJSON(&adminRights); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// فراخوانی لایه سرویس برای به‌روزرسانی حقوق ادمین
	// این تابع فرضی در لایه سرویس، حقوق ادمین را در پایگاه داده به‌روزرسانی می‌کند.
	updatedChat, err := yourService.UpdateChatAdminRights(c.Request.Context(), chatID, userID, &adminRights)
	if err != nil {
		// مدیریت انواع خطاها:
		switch err.Error() {
		case "not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat or user not found"})
		case "permission denied":
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin rights"})
		}
		return
	}

	// ارسال پاسخ موفقیت آمیز به همراه شیء به‌روزرسانی‌شده چت
	c.JSON(http.StatusOK, updatedChat)
}

// UpdateChatBannedRights godoc
// @Summary Update banned rights for a user in a chat
// @Description Updates the banned rights of a user within a chat.
// @Accept  json
// @Produce  json
// @Param chatId path string true "ID of the chat"
// @Param userId path string true "ID of the user whose rights are being updated"
// @Param request body ChatBannedRights true "The banned rights to apply"
// @Success 200 {object} Chat "The updated chat object"
// @Failure 400 {object} map[string]string "Invalid request or parameters"
// @Failure 404 {object} map[string]string "Chat or user not found"
// @Failure 403 {object} map[string]string "Permission denied"
// @Router /chats/{chatId}/users/{userId}/bannedRights [put]
func UpdateChatBannedRights(c *gin.Context) {
	// 1. دریافت شناسه های چت و کاربر از پارامترهای URL
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

	// 2. اتصال بدنه JSON درخواست به ساختار ChatBannedRights
	var bannedRights telegram2.ChatBannedRights
	if err := c.ShouldBindJSON(&bannedRights); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. فراخوانی لایه سرویس برای به‌روزرسانی حقوق ممنوعیت
	updatedChat, err := yourService.UpdateChatBannedRights(c.Request.Context(), chatID, userID, &bannedRights)
	if err != nil {
		// مدیریت انواع خطاها:
		switch err.Error() {
		case "not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat or user not found"})
		case "permission denied":
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update banned rights"})
		}
		return
	}

	// 4. ارسال پاسخ موفقیت‌آمیز به همراه شیء به‌روزرسانی‌شده چت
	c.JSON(http.StatusOK, updatedChat)
}
