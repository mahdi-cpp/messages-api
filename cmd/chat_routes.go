package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
)

func chatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {

	router.POST("/api/chats", chatHandler.Create)

	router.GET("/api/chats", chatHandler.Read)

	//router.GET("/api/chats", chatHandler.Read)
	//router.GET("/api/chats/:chatId/messages/:messageId", chatHandler.ReadChatMessage)
	//router.GET("/api/chats/:chatId/messages/", chatHandler.ReadChatMessages)

	router.PATCH("/api/chats/:chatId", chatHandler.Update)
	router.PATCH("/api/chats/bulk-update", chatHandler.BuckUpdate)

	//router.PATCH("/api/members/:chatId", chatHandler.UpdateMembers)

	router.DELETE("/api/chats/:chatId", chatHandler.Delete)
	router.POST("/api/chats/bulk-delete", chatHandler.BuckDelete)
}
