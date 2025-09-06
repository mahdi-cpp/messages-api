package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

// setupRoutes takes the router and dependencies as arguments
// to ensure all routes are registered on the same instance.
func setupRoutes(router *gin.Engine, appManager *application.Manager, chatHandler *handlers.ChatHandler) {

	// WebSocket route
	router.GET("/ws", func(c *gin.Context) {
		handlers.ServeWs(appManager, c.Writer, c.Request)
	})

	// Static files route
	router.Static("/files/", "./static")

	chatRoutes(router, chatHandler)
}

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

func messageRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {
	router.POST("/api/chats/", chatHandler.Create)
}
