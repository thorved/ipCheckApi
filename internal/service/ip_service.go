package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ipCheckApi/internal/models"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type IPService struct {
	cache        map[string]*models.CachedIPInfo
	cacheMutex   sync.RWMutex
	providers    []models.APIProvider
	currentIndex int
	indexMutex   sync.Mutex
}

// NewIPService creates a new IP service instance
func NewIPService() *IPService {
	return &IPService{
		cache: make(map[string]*models.CachedIPInfo),
		providers: []models.APIProvider{
			{
				Name:    "iplocation.net",
				URL:     "https://www.iplocation.net/get-ipdata",
				Enabled: true,
			},
			// Add more providers here in the future
		},
		currentIndex: 0,
	}
}

// GetIPInfo retrieves IP information with caching
func (s *IPService) GetIPInfo(ip string, ipvType string) (*models.IPInfo, error) {
	// Check cache first
	if cachedInfo := s.getCachedInfo(ip); cachedInfo != nil {
		return &cachedInfo.Data, nil
	}

	// If not in cache or expired, fetch from external API
	ipInfo, err := s.fetchFromProviders(ip, ipvType)
	if err != nil {
		return nil, err
	}

	// Cache the result for 1 hour
	s.cacheInfo(ip, *ipInfo)

	return ipInfo, nil
}

// getCachedInfo retrieves cached IP information if available and not expired
func (s *IPService) getCachedInfo(ip string) *models.CachedIPInfo {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	cachedInfo, exists := s.cache[ip]
	if !exists {
		return nil
	}

	if cachedInfo.IsExpired() {
		// Remove expired entry
		delete(s.cache, ip)
		return nil
	}

	return cachedInfo
}

// cacheInfo stores IP information in cache
func (s *IPService) cacheInfo(ip string, info models.IPInfo) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	now := time.Now()
	s.cache[ip] = &models.CachedIPInfo{
		Data:      info,
		CachedAt:  now,
		ExpiresAt: now.Add(time.Hour), // Cache for 1 hour
	}
}

// fetchFromProviders attempts to fetch IP information from available providers using round-robin
func (s *IPService) fetchFromProviders(ip string, ipvType string) (*models.IPInfo, error) {
	enabledProviders := s.getEnabledProviders()
	if len(enabledProviders) == 0 {
		return nil, fmt.Errorf("no enabled providers available")
	}

	// Set default ipvType if not provided
	if ipvType == "" {
		ipvType = "4"
	}

	// Try each provider starting from the current index (round-robin)
	for i := 0; i < len(enabledProviders); i++ {
		provider := s.getNextProvider(enabledProviders)

		ipInfo, err := s.fetchFromProvider(provider, ip, ipvType)
		if err == nil {
			return ipInfo, nil
		}

		// Log the error but continue to next provider
		fmt.Printf("Provider %s failed: %v\n", provider.Name, err)
	}

	return nil, fmt.Errorf("all providers failed to fetch IP information")
}

// getEnabledProviders returns a list of enabled providers
func (s *IPService) getEnabledProviders() []models.APIProvider {
	var enabled []models.APIProvider
	for _, provider := range s.providers {
		if provider.Enabled {
			enabled = append(enabled, provider)
		}
	}
	return enabled
}

// getNextProvider returns the next provider in round-robin fashion
func (s *IPService) getNextProvider(providers []models.APIProvider) models.APIProvider {
	s.indexMutex.Lock()
	defer s.indexMutex.Unlock()

	provider := providers[s.currentIndex]
	s.currentIndex = (s.currentIndex + 1) % len(providers)
	return provider
}

// fetchFromProvider fetches IP information from a specific provider
func (s *IPService) fetchFromProvider(provider models.APIProvider, ip string, ipvType string) (*models.IPInfo, error) {
	switch provider.Name {
	case "iplocation.net":
		return s.fetchFromIPLocationNet(provider.URL, ip, ipvType)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider.Name)
	}
}

// fetchFromIPLocationNet fetches IP information from iplocation.net API
func (s *IPService) fetchFromIPLocationNet(apiURL, ip, ipvType string) (*models.IPInfo, error) {
	// Prepare form data
	formData := url.Values{}
	formData.Set("ipv", ipvType)
	formData.Set("ip", ip)
	formData.Set("source", "ip2location")

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var apiResponse models.IPLocationNetResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to standardized format
	ipInfo := apiResponse.ToIPInfo()
	return &ipInfo, nil
}

// AddProvider adds a new API provider
func (s *IPService) AddProvider(provider models.APIProvider) {
	s.providers = append(s.providers, provider)
}

// GetProviders returns all configured providers
func (s *IPService) GetProviders() []models.APIProvider {
	return s.providers
}

// EnableProvider enables/disables a provider by name
func (s *IPService) EnableProvider(name string, enabled bool) error {
	for i, provider := range s.providers {
		if provider.Name == name {
			s.providers[i].Enabled = enabled
			return nil
		}
	}
	return fmt.Errorf("provider %s not found", name)
}

// ClearCache clears all cached entries
func (s *IPService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache = make(map[string]*models.CachedIPInfo)
}

// GetCacheStats returns cache statistics
func (s *IPService) GetCacheStats() map[string]interface{} {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	stats := map[string]interface{}{
		"total_entries": len(s.cache),
		"entries":       make([]map[string]interface{}, 0),
	}

	for ip, info := range s.cache {
		entry := map[string]interface{}{
			"ip":         ip,
			"cached_at":  info.CachedAt.Unix(),
			"expires_at": info.ExpiresAt.Unix(),
			"expired":    info.IsExpired(),
		}
		stats["entries"] = append(stats["entries"].([]map[string]interface{}), entry)
	}

	return stats
}
