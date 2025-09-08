package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
)

func chatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {

	router.POST("/api/chats/", chatHandler.Create)

	router.GET("/api/chats/:id", chatHandler.Read)
	router.GET("/api/chats", chatHandler.ReadAll)

	router.PATCH("/api/chats/:id", chatHandler.Update)
	router.PATCH("/api/chats/bulk-update", chatHandler.BuckUpdate)

	//router.PATCH("/api/members/:chatId", chatHandler.UpdateMembers)

	router.DELETE("/api/chats/:id", chatHandler.Delete)
	router.POST("/api/chats/bulk-delete", chatHandler.BuckDelete)
}
