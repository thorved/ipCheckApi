# IP Check API

A robust IP geolocation API service with caching and multiple provider support.

## Features

- **IP Geolocation**: Get detailed information about any IP address including country, city, region, ISP, and coordinates
- **Smart Caching**: 1-hour cache to reduce external API calls for the same IP
- **Round-Robin Providers**: Multiple API provider support with automatic failover
- **RESTful API**: Clean REST endpoints with JSON responses
- **Provider Management**: Enable/disable providers dynamically
- **Cache Management**: View cache statistics and clear cache when needed

## API Endpoints

### IP Lookup

#### POST /api/v1/ip/lookup
Lookup IP information using JSON request body.

**Request Body:**
```json
{
  "ip": "125.99.184.36",
  "ipv_type": "4"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "ipAddress": "125.99.184.36",
    "countryName": "India",
    "countryCode": "IN",
    "regionName": "Delhi",
    "cityName": "Delhi",
    "isp": "Hathway Cable and Datacom Limited",
    "latitude": 28.666774749755859,
    "longitude": 77.216682434082031,
    "timestamp": 1643723400
  }
}
```

#### GET /api/v1/ip/lookup?ip=IP_ADDRESS&ipv_type=4
Lookup IP information using query parameters.

**Parameters:**
- `ip` (required): IP address to lookup
- `ipv_type` (optional): IP version type ("4" or "6"), defaults to "4"

### Cache Management

#### GET /api/v1/cache/stats
Get cache statistics including all cached entries and their expiration status.

#### DELETE /api/v1/cache
Clear all cached entries.

### Provider Management

#### GET /api/v1/providers
Get all configured API providers and their status.

#### PUT /api/v1/providers/enable
Enable or disable a specific provider.

**Request Body:**
```json
{
  "name": "iplocation.net",
  "enabled": true
}
```

### Health Check

#### GET /health or GET /api/v1/health
Check API health status.

## Setup and Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd ipCheckApi
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

4. **Test the API:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/ip/lookup \
     -H "Content-Type: application/json" \
     -d '{"ip": "125.99.184.36", "ipv_type": "4"}'
   ```

## Configuration

### Adding New Providers

The system is designed to support multiple IP geolocation providers. Currently supported:

- **iplocation.net**: Primary provider using IP2Location database

To add new providers, implement the provider-specific logic in the `service/ip_service.go` file and update the `fetchFromProvider` method.

### Cache Settings

- **Cache Duration**: 1 hour (configurable in `service/ip_service.go`)
- **Cache Key**: IP address
- **Cache Storage**: In-memory (can be extended to Redis for distributed caching)

## Architecture

```
cmd/
  server/
    main.go          # Application entry point
internal/
  controller/
    ip_controller.go # HTTP handlers and request validation
  models/
    ip_info.go       # Data structures and models
  routes/
    routes.go        # Route definitions and middleware
  service/
    ip_service.go    # Business logic, caching, and provider management
```

## API Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error Response:**
```json
{
  "error": "Error message",
  "details": "Detailed error information"
}
```

## Caching Strategy

1. **Check Cache**: First check if IP information exists in cache and is not expired
2. **Fetch from Provider**: If not cached or expired, fetch from external providers using round-robin
3. **Update Cache**: Store the result in cache with 1-hour expiration
4. **Failover**: If one provider fails, automatically try the next available provider

## Error Handling

- **Invalid IP Address**: Returns 400 Bad Request
- **Provider Failures**: Automatic failover to next provider
- **All Providers Down**: Returns 500 Internal Server Error
- **Invalid Request Format**: Returns 400 Bad Request with details

## Development

### Project Structure
The project follows clean architecture principles with separation of concerns:

- **Models**: Data structures and business entities
- **Service**: Business logic and external API integration
- **Controller**: HTTP request handling and validation
- **Routes**: URL routing and middleware configuration

### Dependencies
- **Gin**: HTTP web framework
- **Standard Library**: HTTP client, JSON handling, time management

## Future Enhancements

1. **Database Integration**: Store cache in persistent storage (Redis/PostgreSQL)
2. **Rate Limiting**: Implement rate limiting per IP/API key
3. **Authentication**: Add API key authentication
4. **Metrics**: Add Prometheus metrics for monitoring
5. **More Providers**: Add support for additional IP geolocation providers
6. **Geographic Filtering**: Add filtering by country/region
7. **Batch Processing**: Support multiple IP lookups in single request