package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/messages-api/internal/api/handlers"
	"github.com/mahdi-cpp/messages-api/internal/application"
	"github.com/mahdi-cpp/messages-api/internal/collections/chat"
	"github.com/mahdi-cpp/messages-api/internal/config"
)

var port = 50151

func main() {

	config.Init()

	// 1. create a single router instance for the entire application.
	router := gin.Default()
	start := time.Now()

	appManager, err := application.NewApplicationManager()
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("The operation took %s\n", elapsed)

	// Force a garbage collection to get a more accurate picture
	runtime.GC()

	//addChat(appManager, 10000)

	//all, err := appManager.ChatCollectionManager.ReadAll()
	//if err != nil {
	//	return
	//}

	//for _, chat1 := range all {
	//
	//	updateOptions := chat.UpdateOptions{
	//		ChatIDs: []uuid.UUID{
	//			chat1.ID,
	//		},
	//		Title: "ali 09355512617",
	//		Type:  "private",
	//
	//		MembersUpdates: []update.NestedFieldUpdate[chat.Member]{
	//			{
	//				ID: config.Golnar,
	//				Value: func(m *chat.Member) {
	//					m.Role = "admin" // Assuming Member has a Role field
	//				},
	//			},
	//		},
	//	}
	//
	//	err := appManager.UpdateChats(updateOptions)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

	fmt.Println("count: ", appManager.ChatCollectionManager.Count())

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

func addChat(manager *application.AppManager, count int) {

	for i := 0; i < count; i++ {
		chat1 := &chat.Chat{
			Title: fmt.Sprintf("Chat  %d", i),
			Type:  "private",
			Members: []chat.Member{
				{
					UserID: config.Mahdi,
					Role:   "admin",
				},
				{
					UserID: config.Ali,
					Role:   "admin",
				},
				{
					UserID: config.Golnar,
					Role:   "member",
				},
			},
		}

		_, err := manager.ChatCreate(chat1)
		if err != nil {
			return
		}
	}
}
