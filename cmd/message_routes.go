package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
)

func messageRoutes(router *gin.Engine, messageHandler *handlers.MessageHandler) {

	router.POST("/api/messages", messageHandler.Create)
	router.GET("/api/messages", messageHandler.Read)

	//router.GET("/api/messages/:messageId/chats/chatId", messageHandler.Read)
	//router.GET("/api/messages/chats/chatId", messageHandler.ReadAll)

	router.PATCH("/api/messages/:id", messageHandler.Update)
	router.PATCH("/api/messages/bulk-update", messageHandler.BuckUpdate)

	router.DELETE("/api/messages/:id", messageHandler.Delete)
	router.POST("/api/messages/bulk-delete", messageHandler.BuckDelete)
}
