# API Gateway

A simple, high-performance API Gateway written in Go that routes requests to backend microservices.

## Features

- Reverse proxy to multiple backend services
- Rate limiting per client IP
- CORS support
- Request logging
- Graceful shutdown
- Health check endpoint
- JSON-based configuration

## Project Structure

```
.
├── cmd/
│   └── apigateway/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── gateway/
│   │   └── gateway.go        # Core gateway logic
│   └── middleware/
│       └── middleware.go     # HTTP middleware (logging, CORS, rate limiting)
├── pkg/
│   └── router/
│       └── router.go         # HTTP router implementation
├── config.json               # Configuration file
├── go.mod                    # Go module definition
└── README.md
```

## Prerequisites

- Go 1.21 or higher

## Installation

1. Install Go from [golang.org](https://golang.org/dl/)

2. Initialize the Go module and install dependencies:

```bash
go mod init apigateway
go mod tidy
```

## Configuration

The gateway is configured via `config.json`. Example configuration:

```json
{
  "server": {
    "address": ":8080",
    "read_timeout": 15000000000,
    "write_timeout": 15000000000,
    "idle_timeout": 60000000000
  },
  "services": [
    {
      "name": "user-service",
      "path": "/api/users/*",
      "target": "http://localhost:9001"
    },
    {
      "name": "product-service",
      "path": "/api/products/*",
      "target": "http://localhost:9002"
    }
  ],
  "rate_limit": 100
}
```

Configuration options:

- `server.address`: The address the gateway listens on (default: `:8080`)
- `server.read_timeout`: Maximum duration for reading the entire request (in nanoseconds)
- `server.write_timeout`: Maximum duration for writing the response (in nanoseconds)
- `server.idle_timeout`: Maximum idle time between requests (in nanoseconds)
- `services`: Array of backend services to proxy
  - `name`: Service identifier
  - `path`: URL path pattern (use `/*` for wildcard matching)
  - `target`: Backend service URL
- `rate_limit`: Maximum requests per second per IP address

## Running the Gateway

```bash
go run cmd/apigateway/main.go
```

Or build and run:

```bash
go build -o apigateway cmd/apigateway/main.go
./apigateway
```

## Endpoints

- `GET /health` - Health check endpoint
- All configured service routes are proxied according to the configuration

## Example Usage

Once the gateway is running on port 8080:

```bash
# Health check
curl http://localhost:8080/health

# Access proxied services
curl http://localhost:8080/api/users/123
curl http://localhost:8080/api/products/456
```

## Development

To add a new middleware:

1. Create your middleware function in `internal/middleware/middleware.go`
2. Register it in `internal/gateway/gateway.go` using `gw.router.Use(middleware.YourMiddleware)`

To add new routes:

1. Update `config.json` with the new service configuration
2. Restart the gateway

## Dependencies

The project uses minimal external dependencies:

- `golang.org/x/time/rate` - For rate limiting functionality

## Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o apigateway cmd/apigateway/main.go
```

## Docker Support (Coming Soon)

A Dockerfile will be added for containerized deployment.

## License

MIT
