package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
	"github.com/mahdi-cpp/messages-api/internal/application"
)

func main() {

	appManager, err := application.NewApplicationManager()
	if err != nil {
		log.Fatal(err)
	}

	chats, err := appManager.GetUserChats("018f3a8b-1b32-7293-c1d4-8765f4d1e2f3")
	if err != nil {
		return
	}

	if len(chats) == 0 {
		log.Printf("No chat found")
	} else {
		fmt.Println("-----------------------")
		for _, chat := range chats {
			fmt.Println(chat.Title)
		}
		fmt.Println("-----------------------")
	}

	ginInit()
	chatHandler := handlers.NewChatHandler(appManager)
	chatRouter(chatHandler)

	// Convert http.HandleFunc to a Gin handler
	// Gin's handlers have the signature func(*gin.Context)
	router.GET("/ws", func(c *gin.Context) {
		// Get the http.ResponseWriter and *http.Request from the Gin context
		w := c.Writer
		r := c.Request

		// Call your original handler function
		handlers.ServeWs(appManager, w, r)
	})

	// Convert http.FileServer to Gin's router.Static
	router.Static("/files/", "./static")

	startServer(router)
}

func webRouter(chatHandler *handlers.ChatHandler) *gin.Engine {

	api := router.Group("")

	api.GET("/", chatHandler.Create)
	api.POST("update", chatHandler.Update)
	api.POST("delete", chatHandler.Delete)
	//api.POST("getFilters", chatHandler.GetFilter)

	return router
}

func chatRouter(chatHandler *handlers.ChatHandler) *gin.Engine {

	api := router.Group("/api/")

	api.PUT("chats/{chatId}/users/{userId}/createChats", chatHandler.Create)
	api.POST("update", chatHandler.Update)
	api.POST("delete", chatHandler.Delete)
	//api.POST("getFilters", chatHandler.GetFilter)

	return router
}
