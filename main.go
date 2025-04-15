package main

import (
	"api/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", handlers.HealthCheck)

	// Start server
	r.Run(":8080")
}
