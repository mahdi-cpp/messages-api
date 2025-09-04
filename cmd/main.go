package main

import (
	"fmt"
	"log"
	"net/http"

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

	// Set up HTTP handlers
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(appManager, w, r)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Start server
	log.Println("Server starting on :8089")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("Server error: ", err)
	}

	ginInit()
	chatHandler := handlers.NewChatHandler(appManager)
	chatRouter(chatHandler)

	startServer(router)
}

func chatRouter(chatHandler *handlers.ChatHandler) *gin.Engine {

	api := router.Group("/api/chat/")

	api.POST("create", chatHandler.Create)
	api.POST("update", chatHandler.Update)
	api.POST("delete", chatHandler.Delete)
	api.POST("getFilters", chatHandler.GetFilter)

	return router
}
