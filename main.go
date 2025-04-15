package main

import (
	"root/api/handlers"
	"root/api/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Initialize repository and handler
	chatroomRepo := repositories.NewMemoryChatroomRepository()
	chatroomHandler := handlers.NewChatroomHandler(chatroomRepo)

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
	}

	// Start server
	r.Run(":8080")
}
