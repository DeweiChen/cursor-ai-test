package main

import (
	"root/api/handlers"
	"root/api/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize repositories
	chatroomRepo := repositories.NewMemoryChatroomRepository()
	messageRepo := repositories.NewMemoryMessageRepository()

	// Initialize handlers
	chatroomHandler := handlers.NewChatroomHandler(chatroomRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo, chatroomRepo)

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Chatroom routes
	chatroomGroup := r.Group("/api/chatrooms")
	{
		chatroomGroup.POST("", chatroomHandler.CreateChatroom)
		chatroomGroup.GET("", chatroomHandler.GetAllChatrooms)
		chatroomGroup.GET("/:id", chatroomHandler.GetChatroom)
		chatroomGroup.PUT("/:id", chatroomHandler.UpdateChatroom)
		chatroomGroup.DELETE("/:id", chatroomHandler.DeleteChatroom)

		// Message routes
		chatroomGroup.POST("/:id/messages", messageHandler.CreateMessage)
		chatroomGroup.GET("/:id/messages", messageHandler.GetMessagesByChatroomID)
	}

	// Start server
	r.Run(":8080")
}
