package routes

import (
	"ipCheckApi/internal/controller"
	"ipCheckApi/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, ipService *service.IPService) {
	ipController := controller.NewIPController(ipService)

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// IP information endpoints
		v1.POST("/ip/lookup", ipController.GetIPInfo)
		v1.GET("/ip/lookup", ipController.GetIPInfoByQuery)

		// Cache management endpoints
		v1.GET("/cache/stats", ipController.GetCacheStats)
		v1.DELETE("/cache", ipController.ClearCache)

		// Provider management endpoints
		v1.GET("/providers", ipController.GetProviders)
		v1.PUT("/providers/enable", ipController.EnableProvider)

		// Health check
		v1.GET("/health", ipController.Health)
	}

	// Root level health check
	router.GET("/health", ipController.Health)

	// Add CORS middleware for cross-origin requests
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}
