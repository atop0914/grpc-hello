# gRPC-Hello

A production-ready gRPC microservice with HTTP/JSON gateway in Go, featuring multi-language support, statistics tracking, and comprehensive monitoring.

## 🚀 Features

- **gRPC Service**: High-performance RPC communication
- **HTTP/JSON Gateway**: Automatic RESTful API via gRPC-Gateway
- **Multi-language Support**: International greetings in 9+ languages
- **Real-time Statistics**: Request counting and analytics
- **Health Monitoring**: Built-in health checks
- **Prometheus Metrics**: Production-grade observability
- **Graceful Shutdown**: Safe service termination
- **Docker Ready**: Container-first design
- **Cross-platform**: Build for Linux/macOS/Windows

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Client   │────│  gRPC-Gateway    │────│   gRPC Server   │
│   (REST/JSON)   │    │ (HTTP ↔ gRPC)    │    │  (Protocol     │
└─────────────────┘    │   Translation    │    │   Buffers)     │
                       └──────────────────┘    └─────────────────┘
                              │                        │
                       ┌──────────────────┐    ┌─────────────────┐
                       │     Gin Router   │    │  Stats &        │
                       │  (middleware,    │    │   Monitoring    │
                       │   metrics)       │    │   Engine        │
                       └──────────────────┘    └─────────────────┘
```

## 🛠️ Prerequisites

- Go 1.22+
- Git

## 🚀 Quick Start

### Clone and Build

```bash
git clone git@github.com:atop0914/grpc-hello.git
cd grpc-hello

# Install dependencies
make deps

# Build the service
make build

# Run the service
make run
```

### Default Endpoints

- **gRPC**: `localhost:8080`
- **HTTP**: `localhost:8090`
- **Metrics**: `localhost:8090/metrics`
- **Health**: `localhost:8090/health`

## 📡 API Usage

### gRPC Client

```bash
# Basic greeting
go run client/client.go

# Custom name
go run client/client.go --name="Alice"

# Custom server
go run client/client.go --addr="localhost:9090" --name="Bob"
```

### HTTP API

```bash
# Basic greeting (English)
curl -X POST http://localhost:8090/rpc/v1/sayHello \
  -H "Content-Type: application/json" \
  -d '{"name_test": "World", "language": "en"}'

# International greeting (Chinese)
curl -X POST http://localhost:8090/rpc/v1/sayHello \
  -H "Content-Type: application/json" \
  -d '{"name_test": "世界", "language": "zh"}'

# Multiple greetings
curl -X POST http://localhost:8090/rpc/v1/sayHelloMultiple \
  -H "Content-Type: application/json" \
  -d '{"names": ["Alice", "Bob"], "common_message": "Welcome!"}'

# Get statistics
curl -X GET http://localhost:8090/rpc/v1/greetingStats
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GRPC_PORT` | gRPC server port | `8080` |
| `HTTP_PORT` | HTTP server port | `8090` |
| `ENABLE_DEBUG` | Debug mode | `false` |
| `SERVER_TIMEOUT` | Server timeout (seconds) | `30` |
| `LOG_LEVEL` | Logging level | `info` |
| `ENABLE_REFLECTION` | gRPC reflection | `false` |
| `ENABLE_STATS` | Statistics tracking | `true` |

### Runtime Configuration

```bash
# Custom ports
GRPC_PORT=9090 HTTP_PORT=9091 make run

# Enable debug mode
ENABLE_DEBUG=true make run
```

## 🏗️ Building

### Single Platform

```bash
# Build for current platform
make build

# Run directly
go run main.go
```

### Cross-platform Builds

```bash
# Build for all platforms
make build-all

# Build specific platforms
make build-linux    # Linux binary
make build-mac      # macOS binary  
make build-windows  # Windows binary

# Clean artifacts
make clean
```

## 🐳 Docker Deployment

```bash
# Build image
docker build -t grpc-hello .

# Run container
docker run -p 8080:8080 -p 8090:8090 grpc-hello

# Run with custom configuration
docker run -e GRPC_PORT=9090 -e HTTP_PORT=9091 -p 9090:9090 -p 9091:9091 grpc-hello
```

## 🌍 Supported Languages

| Code | Language | Greeting |
|------|----------|----------|
| `en` | English | Hello |
| `zh` | Chinese | 你好 |
| `es` | Spanish | Hola |
| `fr` | French | Bonjour |
| `ja` | Japanese | こんにちは |
| `ko` | Korean | 안녕하세요 |
| `ru` | Russian | Привет |
| `de` | German | Hallo |
| `it` | Italian | Ciao |

## 📊 Monitoring & Observability

### Metrics

- **Prometheus Endpoint**: `GET /metrics`
- **Health Check**: `GET /health`
- **Statistics**: `GET /rpc/v1/greetingStats`

### Key Metrics Tracked

- Total requests served
- Language distribution
- Request patterns
- Service health status

## 📁 Project Structure

```
grpc-hello/
├── main.go               # Core gRPC service
├── config/               # Configuration management
│   └── config.go
├── client/               # gRPC client example
│   └── client.go
├── proto/                # Protocol buffers
│   └── helloworld/
│       ├── hello_world.proto      # Service definition
│       ├── hello_world.pb.go      # Generated Go
│       ├── hello_world_grpc.pb.go # Generated gRPC
│       └── hello_world.pb.gw.go   # Generated gateway
├── route/                # HTTP routes
│   └── route.go
├── Makefile              # Build automation
├── Dockerfile            # Container spec
├── go.mod                # Dependencies
└── README.md
```

## 🧪 Testing

```bash
# Run tests
make test

# Install dependencies
make deps
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.