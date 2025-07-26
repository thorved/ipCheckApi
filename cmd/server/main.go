package main

import (
	"ipCheckApi/internal/routes"
	"ipCheckApi/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize services
	ipService := service.NewIPService()

	// Setup Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, ipService)

	// Start server
	port := ":8080"
	log.Printf("Starting IP Check API server on port %s", port)
	log.Printf("Health check available at: http://localhost%s/health", port)
	log.Printf("API documentation available at: http://localhost%s/api/v1/health", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
