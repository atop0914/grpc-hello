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
- ✅ **Cross-platform Builds**: Easy deployment across systems

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

- Go 1.22+
- Protocol Buffers (protobuf) compiler (optional, for development)
- Git

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

## Quick Start

### Running the Server

```bash
# Run with default settings
make run

# Or run the built binary
./grpc-hello

# With environment variables
GRPC_PORT=9090 HTTP_PORT=9091 ENABLE_DEBUG=true make run
```

### Using the Client

```bash
# Run the client
go run client/client.go

# With custom name
go run client/client.go --name="Alice"

# With custom server address
go run client/client.go --addr="localhost:8090" --name="Bob"
```

### Available Endpoints

#### gRPC Endpoints
- `localhost:8080` - Main gRPC service

#### HTTP Endpoints
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics
- `POST /rpc/v1/sayHello` - Basic greeting endpoint
- `POST /rpc/v1/sayHelloMultiple` - Multiple greetings endpoint
- `GET /rpc/v1/greetingStats` - Statistics endpoint

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
- `LOG_LEVEL`: Log level (default: info)
- `ENABLE_REFLECTION`: Enable gRPC reflection (default: matches debug mode)
- `ENABLE_STATS`: Enable statistics tracking (default: true)

## Building and Deployment

### Build Options

```bash
# Build for current platform
make build

# Build for Linux
make build-linux

# Build for macOS
make build-mac

# Build for Windows
make build-windows

# Build all platforms
make build-all

# Clean build artifacts
make clean
```

### Docker Deployment

```bash
# Build Docker image
docker build -t grpc-hello .

# Run in Docker
docker run -p 8080:8080 -p 8090:8090 grpc-hello
```

### Cross-platform Binary Distribution

```bash
# Build binaries for all platforms
make build-all

# Binaries will be created as:
# - grpc-hello-linux
# - grpc-hello-mac
# - grpc-hello-windows.exe
```

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
├── proto/
│   └── helloworld/
│       └── hello_world.proto # Protocol buffer definitions
│       └── hello_world.pb.go # Generated Go code
│       └── hello_world_grpc.pb.go # Generated gRPC code
│       └── hello_world.pb.gw.go # Generated gateway code
├── route/
│   └── route.go              # HTTP route definitions
└── README.md
```

## Development

### Building and Testing

```bash
# Build the application
make build

# Run tests
make test

# Install dependencies
make deps
```

## Internationalization

Supported languages:
- English: "Hello"
- Chinese: "你好"
- Spanish: "Hola"
- French: "Bonjour"
- Japanese: "こんにちは"
- Korean: "안녕하세요"
- Russian: "Привет"
- German: "Hallo"
- Italian: "Ciao"

## Monitoring

The service exposes Prometheus metrics at `/metrics` endpoint.

## License

MIT License - See LICENSE file for details.