# gRPC-Hello

A comprehensive gRPC demonstration project in Go with gRPC-Gateway integration, providing both gRPC and RESTful HTTP APIs. This enhanced version includes advanced features like streaming, statistics, and internationalization.

## Features

- ✅ **gRPC Service**: High-performance RPC communication
- ✅ **gRPC-Gateway**: HTTP/JSON RESTful API gateway
- ✅ **Multiple Endpoints**: Both gRPC and HTTP interfaces
- ✅ **Health Checks**: Built-in health monitoring
- ✅ **Metrics**: Prometheus metrics integration
- ✅ **Graceful Shutdown**: Proper cleanup on termination
- ✅ **Configuration**: Flexible configuration options via env vars
- ✅ **Error Handling**: Comprehensive error management
- ✅ **Streaming Support**: Server-side and client-side streaming
- ✅ **Statistics Tracking**: Request counting and analysis
- ✅ **Internationalization**: Multi-language greetings
- ✅ **Rate Limiting**: Prevents abuse
- ✅ **API Documentation**: Built-in API docs endpoint
- ✅ **Docker Support**: Ready for container deployment

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Client   │────│  gRPC-Gateway    │────│   gRPC Server   │
│   (REST/JSON)   │    │ (HTTP/JSON ->    │    │  (Protocol     │
└─────────────────┘    │  gRPC Protocol)  │    │   Buffers)     │
                       └──────────────────┘    └─────────────────┘
                              │                        │
                       ┌──────────────────┐    ┌─────────────────┐
                       │     Gin Router   │    │  Statistics &   │
                       │  (middleware,    │    │   Monitoring    │
                       │   metrics, etc)  │    │   Engine        │
                       └──────────────────┘    └─────────────────┘
```

## Prerequisites

- Go 1.23+
- Protocol Buffers (protobuf) compiler
- Buf CLI (optional, for protobuf generation)

## Installation

```bash
# Clone the repository
git clone https://github.com/atop0914/grpc-hello.git
cd grpc-hello

# Install dependencies
make deps

# Build the application
make build

# Or run directly
go run main.go
```

## Usage

### Running the Server

```bash
# Run with default settings
make run

# Or run the built binary
./grpc-hello

# With environment variables
GRPC_PORT=9090 HTTP_PORT=9091 ENABLE_DEBUG=true make run

# With custom configuration
go run main.go --grpc-port=9090 --http-port=9091 --debug
```

### Available Endpoints

#### gRPC Endpoints
- `localhost:8080` - Main gRPC service
- Supports unary, server streaming, client streaming, and bidirectional streaming

#### HTTP Endpoints
- `GET /` - Welcome page with API documentation
- `GET /health` - Basic health check
- `GET /healthz` - Detailed health check
- `GET /status` - Service status and memory stats
- `GET /version` - Service version info
- `GET /metrics` - Prometheus metrics
- `GET /docs` - API documentation
- `GET /ping` - Ping endpoint
- `POST /echo` - Echo endpoint for testing
- `POST /rpc/v1/sayHello` - Basic greeting endpoint
- `POST /rpc/v1/sayHelloMultiple` - Multiple greetings endpoint
- `GET /rpc/v1/greetingStats` - Statistics endpoint

### Using the Client

```bash
# Run the client
go run client/client.go

# With custom name
go run client/client.go --name="Alice"

# With custom server address
go run client/client.go --addr="localhost:9090" --name="Bob"
```

### Example API Calls

```bash
# Basic greeting
curl -X POST http://localhost:8090/rpc/v1/sayHello \
  -H "Content-Type: application/json" \
  -d '{"name_test": "World", "language": "en"}'

# Multiple greetings
curl -X POST http://localhost:8090/rpc/v1/sayHelloMultiple \
  -H "Content-Type: application/json" \
  -d '{"names": ["Alice", "Bob"], "common_message": "Welcome!"}'

# Get statistics
curl -X GET http://localhost:8090/rpc/v1/greetingStats

# International greeting
curl -X POST http://localhost:8090/rpc/v1/sayHello \
  -H "Content-Type: application/json" \
  -d '{"name_test": "世界", "language": "zh"}'
```

## Configuration Options

The application supports the following configuration options:

### Environment Variables
- `GRPC_PORT`: Port for gRPC server (default: 8080)
- `HTTP_PORT`: Port for HTTP server (default: 8090)
- `ENABLE_DEBUG`: Enable debug mode (default: false)
- `SERVER_TIMEOUT`: Server timeout in seconds (default: 30)
- `MAX_CONNECTIONS`: Maximum concurrent connections (default: 1000)
- `LOG_LEVEL`: Log level (default: info)
- `METRICS_ENABLED`: Enable metrics (default: true)
- `METRICS_PATH`: Metrics endpoint path (default: /metrics)
- `ENABLE_REFLECTION`: Enable gRPC reflection (default: matches debug mode)
- `ENABLE_STATS`: Enable statistics tracking (default: true)
- `MAX_GREETINGS`: Maximum greetings per request (default: 100)

### Command Line Flags
- `--grpc-port`: Port for gRPC server (default: 8080)
- `--http-port`: Port for HTTP server (default: 8090)
- `--debug`: Enable debug mode (default: false)

## Project Structure

```
grpc-hello/
├── main.go                   # Main application entry point
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── Makefile                  # Build automation
├── Dockerfile                # Container configuration
├── config/
│   └── config.go             # Configuration management
├── client/
│   └── client.go             # gRPC client example
├── middleware/
│   └── logging.go            # Middleware utilities
├── proto/
│   └── helloworld/
│       └── hello_world.proto # Protocol buffer definitions
├── route/
│   └── route.go              # HTTP route definitions
├── buf.yaml                  # Buf configuration
├── buf.gen.yaml              # Buf generation configuration
└── README.md
```

## Development

### Regenerating Protocol Buffer Files

If you modify the `.proto` files, regenerate the Go code:

```bash
# Using make
make proto-gen

# Or manually
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

protoc --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=. \
  --grpc-gateway_opt=paths=source_relative \
  proto/helloworld/hello_world.proto
```

### Building and Testing

```bash
# Build the application
make build

# Run tests
make test

# Generate protobuf and build
make gen-and-build

# Run in development mode
make run-dev
```

### Adding New Services

1. Define your service in a `.proto` file
2. Generate Go code using `make proto-gen`
3. Implement the server methods in `main.go`
4. Register the service with the gRPC server
5. Register the gateway handler with the gRPC-Gateway

## Monitoring

The service exposes Prometheus metrics at `/metrics` endpoint. The following metrics are available:

- `go_*`: Go runtime metrics
- `promhttp_*`: HTTP server metrics
- Custom statistics tracking greetings and usage patterns

## Performance

- Uses `endless` for zero-downtime deployments
- Optimized JSON marshaling with protojson
- Efficient gRPC communication
- Connection pooling and reuse
- Rate limiting to prevent abuse
- Statistics tracking without performance impact

## Deployment

### Docker

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run

# Or manually:
docker build -t grpc-hello .
docker run -p 8080:8080 -p 8090:8090 grpc-hello
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Ensure all tests pass
6. Submit a pull request

## License

MIT License - See LICENSE file for details.