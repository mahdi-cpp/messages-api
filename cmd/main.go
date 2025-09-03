package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handler"
	"github.com/mahdi-cpp/messages-api/internal/application"
	"github.com/mahdi-cpp/messages-api/internal/hub"
	"github.com/mahdi-cpp/messages-api/internal/storage"
)

func main() {

	appManager, err := application.NewApplicationManager()
	if err != nil {
		log.Fatal(err)
	}
	ginInit()
	chatHandler := handler.NewChatHandler(appManager)
	chatRouter(chatHandler)

	// Initialize hub
	chatHub := hub.NewHub()
	go chatHub.Run()

	// Load all chat metadata on startup
	chats, err := storage.LoadAllChats()
	if err != nil {
		log.Fatalf("Failed to load chats: %v", err)
	}
	log.Printf("Loaded %d chats", len(chats))

	// Set up HTTP handlers
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeWs(chatHub, w, r)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Start server
	log.Println("Server starting on :8089")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatal("Server error: ", err)
	}

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
