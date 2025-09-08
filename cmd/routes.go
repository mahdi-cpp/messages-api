package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

// setupRoutes takes the router and dependencies as arguments
// to ensure all routes are registered on the same instance.
func setupRoutes(router *gin.Engine, appManager *application.Manager,
	chatHandler *handlers.ChatHandler,
	messageHandler *handlers.MessageHandler,
) {

	// WebSocket route
	router.GET("/ws", func(c *gin.Context) {
		handlers.ServeWs(appManager, c.Writer, c.Request)
	})

	// Static files route
	router.Static("/files/", "./static")

	chatRoutes(router, chatHandler)
	messageRoutes(router, messageHandler)
}
