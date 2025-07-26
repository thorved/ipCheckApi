package models

import "time"

// IPInfo represents the standardized response from our API
type IPInfo struct {
	IPAddress   string  `json:"ipAddress"`
	CountryName string  `json:"countryName"`
	CountryCode string  `json:"countryCode"`
	RegionName  string  `json:"regionName"`
	CityName    string  `json:"cityName"`
	ISP         string  `json:"isp"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timestamp   int64   `json:"timestamp"`
}

// CachedIPInfo represents cached IP information with expiry
type CachedIPInfo struct {
	Data      IPInfo    `json:"data"`
	CachedAt  time.Time `json:"cachedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// IsExpired checks if the cached data has expired
func (c *CachedIPInfo) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// IPLocationNetResponse represents the response from iplocation.net API
type IPLocationNetResponse struct {
	IsProxy bool   `json:"isProxy"`
	Source  string `json:"source"`
	Res     struct {
		IPNumber    string  `json:"ipNumber"`
		IPVersion   int     `json:"ipVersion"`
		IPAddress   string  `json:"ipAddress"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		CountryName string  `json:"countryName"`
		CountryCode string  `json:"countryCode"`
		ISP         string  `json:"isp"`
		CityName    string  `json:"cityName"`
		RegionName  string  `json:"regionName"`
	} `json:"res"`
}

// ToIPInfo converts IPLocationNetResponse to standardized IPInfo
func (r *IPLocationNetResponse) ToIPInfo() IPInfo {
	return IPInfo{
		IPAddress:   r.Res.IPAddress,
		CountryName: r.Res.CountryName,
		CountryCode: r.Res.CountryCode,
		RegionName:  r.Res.RegionName,
		CityName:    r.Res.CityName,
		ISP:         r.Res.ISP,
		Latitude:    r.Res.Latitude,
		Longitude:   r.Res.Longitude,
		Timestamp:   time.Now().Unix(),
	}
}

// APIProvider represents an external IP information provider
type APIProvider struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

// IPRequest represents the request structure for IP lookup
type IPRequest struct {
	IP      string `json:"ip" binding:"required"`
	IPVType string `json:"ipv_type,omitempty"` // "4" or "6", defaults to "4"
}
