package controller

import (
	"ipCheckApi/internal/models"
	"ipCheckApi/internal/service"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IPController struct {
	ipService *service.IPService
}

// NewIPController creates a new IP controller
func NewIPController(ipService *service.IPService) *IPController {
	return &IPController{
		ipService: ipService,
	}
}

// GetIPInfo handles IP information lookup requests
func (c *IPController) GetIPInfo(ctx *gin.Context) {
	var req models.IPRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate IP address
	if !isValidIP(req.IP) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid IP address format",
		})
		return
	}

	// Set default IPv type if not provided
	if req.IPVType == "" {
		req.IPVType = "4"
	}

	// Validate IPv type
	if req.IPVType != "4" && req.IPVType != "6" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "IPV type must be '4' or '6'",
		})
		return
	}

	// Get IP information
	ipInfo, err := c.ipService.GetIPInfo(req.IP, req.IPVType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve IP information",
			"details": err.Error(),
		})
		return
	}

	// Return successful response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ipInfo,
	})
}

// GetIPInfoByQuery handles IP information lookup via query parameters
func (c *IPController) GetIPInfoByQuery(ctx *gin.Context) {
	ip := ctx.Query("ip")
	if ip == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "IP parameter is required",
		})
		return
	}

	// Validate IP address
	if !isValidIP(ip) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid IP address format",
		})
		return
	}

	ipvType := ctx.DefaultQuery("ipv_type", "4")

	// Validate IPv type
	if ipvType != "4" && ipvType != "6" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "IPV type must be '4' or '6'",
		})
		return
	}

	// Get IP information
	ipInfo, err := c.ipService.GetIPInfo(ip, ipvType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve IP information",
			"details": err.Error(),
		})
		return
	}

	// Return successful response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ipInfo,
	})
}

// GetCacheStats returns cache statistics
func (c *IPController) GetCacheStats(ctx *gin.Context) {
	stats := c.ipService.GetCacheStats()
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ClearCache clears all cached entries
func (c *IPController) ClearCache(ctx *gin.Context) {
	c.ipService.ClearCache()
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cache cleared successfully",
	})
}

// GetProviders returns all configured providers
func (c *IPController) GetProviders(ctx *gin.Context) {
	providers := c.ipService.GetProviders()
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    providers,
	})
}

// EnableProvider enables or disables a provider
func (c *IPController) EnableProvider(ctx *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Enabled bool   `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	err := c.ipService.EnableProvider(req.Name, req.Enabled)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Provider not found",
			"details": err.Error(),
		})
		return
	}

	action := "disabled"
	if req.Enabled {
		action = "enabled"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Provider " + req.Name + " " + action + " successfully",
	})
}

// Health check endpoint
func (c *IPController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "IP Check API",
		"version": "1.0.0",
	})
}

// isValidIP validates if the provided string is a valid IP address
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}
