package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
	"github.com/mahdi-cpp/messages-api/internal/application"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

var port = 50151

func main() {
	// 1. Create a single router instance for the entire application.
	router := gin.Default()

	appManager, err := application.NewApplicationManager()
	if err != nil {
		log.Fatal(err)
	}

	chats, err := appManager.ReadUserChats(config.Mahdi)
	if err != nil {
		log.Printf("Error getting chats: %v", err)
		// The server will still start, but log the error.
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

	// 2. Instantiate handlers.
	chatHandler := handlers.NewChatHandler(appManager)
	messageHandler := handlers.NewMessageHandler(appManager)

	// 3. Set up all routes using the single router instance.
	setupRoutes(router, appManager,
		chatHandler,
		messageHandler,
	)

	// 4. Start the server with the fully configured router.
	startServer(router)
}

func startServer(router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", port),
		Handler: router,
	}

	// Graceful shutdown logic...
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
