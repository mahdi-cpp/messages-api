package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handler"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

func main() {

	appManager, err := application.NewApplicationManager()
	if err != nil {
		log.Fatal(err)
	}
	ginInit()
	chatHandler := handler.NewChatHandler(appManager)
	chatRouter(chatHandler)

	startServer(router)
}

func chatRouter(chatHandler *handler.ChatHandler) *gin.Engine {

	api := router.Group("/api/chat/")

	api.POST("create", chatHandler.Create)
	api.POST("update", chatHandler.Update)
	api.POST("delete", chatHandler.Delete)
	api.POST("getFilters", chatHandler.GetFilter)

	return router
}
