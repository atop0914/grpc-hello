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

---

# gRPC-Hello

一款生产就绪的gRPC微服务，带有HTTP/JSON网关，支持多语言、统计跟踪和全面监控。

## 🚀 特性

- **gRPC服务**: 高性能RPC通信
- **HTTP/JSON网关**: 通过gRPC-Gateway自动提供RESTful API
- **多语言支持**: 支持9种以上语言的国际问候
- **实时统计**: 请求计数和分析
- **健康监控**: 内置健康检查
- **Prometheus指标**: 生产级可观测性
- **优雅关闭**: 安全的服务终止
- **Docker就绪**: 容器优先设计
- **跨平台**: 支持Linux/macOS/Windows构建

## 🏗️ 架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP客户端    │────│  gRPC-Gateway    │────│   gRPC服务器    │
│   (REST/JSON)   │    │ (HTTP ↔ gRPC)    │    │  (协议缓冲区)   │
└─────────────────┘    │   转换层          │    │                │
                       └──────────────────┘    └─────────────────┘
                              │                        │
                       ┌──────────────────┐    ┌─────────────────┐
                       │     Gin路由      │    │  统计与          │
                       │  (中间件,        │    │   监控引擎       │
                       │   指标)          │    │                │
                       └──────────────────┘    └─────────────────┘
```

## 🛠️ 前置条件

- Go 1.22+
- Git

## 🚀 快速开始

### 克隆并构建

```bash
git clone git@github.com:atop0914/grpc-hello.git
cd grpc-hello

# 安装依赖
make deps

# 构建服务
make build

# 运行服务
make run
```

### 默认端点

- **gRPC**: `localhost:8080`
- **HTTP**: `localhost:8090`
- **指标**: `localhost:8090/metrics`
- **健康检查**: `localhost:8090/health`

## 📡 API使用

### gRPC客户端

```bash
# 基础问候
go run client/client.go

# 自定义名称
go run client/client.go --name="Alice"

# 自定义服务器
go run client/client.go --addr="localhost:9090" --name="Bob"
```

### HTTP API

```bash
# 基础问候 (英文)
curl -X POST http://localhost:8090/rpc/v1/sayHello \
  -H "Content-Type: application/json" \
  -d '{"name_test": "World", "language": "en"}'

# 国际问候 (中文)
curl -X POST http://localhost:8090/rpc/v1/sayHello \
  -H "Content-Type: application/json" \
  -d '{"name_test": "世界", "language": "zh"}'

# 多个问候
curl -X POST http://localhost:8090/rpc/v1/sayHelloMultiple \
  -H "Content-Type: application/json" \
  -d '{"names": ["Alice", "Bob"], "common_message": "Welcome!"}'

# 获取统计信息
curl -X GET http://localhost:8090/rpc/v1/greetingStats
```

## ⚙️ 配置

### 环境变量

| 变量 | 描述 | 默认值 |
|------|------|--------|
| `GRPC_PORT` | gRPC服务器端口 | `8080` |
| `HTTP_PORT` | HTTP服务器端口 | `8090` |
| `ENABLE_DEBUG` | 调试模式 | `false` |
| `SERVER_TIMEOUT` | 服务器超时(秒) | `30` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `ENABLE_REFLECTION` | gRPC反射 | `false` |
| `ENABLE_STATS` | 统计跟踪 | `true` |

### 运行时配置

```bash
# 自定义端口
GRPC_PORT=9090 HTTP_PORT=9091 make run

# 启用调试模式
ENABLE_DEBUG=true make run
```

## 🏗️ 构建

### 单平台构建

```bash
# 为当前平台构建
make build

# 直接运行
go run main.go
```

### 跨平台构建

```bash
# 为所有平台构建
make build-all

# 构建特定平台
make build-linux    # Linux二进制文件
make build-mac      # macOS二进制文件
make build-windows  # Windows二进制文件

# 清理构建产物
make clean
```

## 🐳 Docker部署

```bash
# 构建镜像
docker build -t grpc-hello .

# 运行容器
docker run -p 8080:8080 -p 8090:8090 grpc-hello

# 运行自定义配置
docker run -e GRPC_PORT=9090 -e HTTP_PORT=9091 -p 9090:9090 -p 9091:9091 grpc-hello
```

## 🌍 支持的语言

| 代码 | 语言 | 问候语 |
|------|------|--------|
| `en` | 英语 | Hello |
| `zh` | 中文 | 你好 |
| `es` | 西班牙语 | Hola |
| `fr` | 法语 | Bonjour |
| `ja` | 日语 | こんにちは |
| `ko` | 韩语 | 안녕하세요 |
| `ru` | 俄语 | Привет |
| `de` | 德语 | Hallo |
| `it` | 意大利语 | Ciao |

## 📊 监控与可观测性

### 指标

- **Prometheus端点**: `GET /metrics`
- **健康检查**: `GET /health`
- **统计信息**: `GET /rpc/v1/greetingStats`

### 跟踪的关键指标

- 服务总请求数
- 语言分布
- 请求模式
- 服务健康状态

## 📁 项目结构

```
grpc-hello/
├── main.go               # 核心gRPC服务
├── config/               # 配置管理
│   └── config.go
├── client/               # gRPC客户端示例
│   └── client.go
├── proto/                # 协议缓冲区
│   └── helloworld/
│       ├── hello_world.proto      # 服务定义
│       ├── hello_world.pb.go      # 生成的Go代码
│       ├── hello_world_grpc.pb.go # 生成的gRPC代码
│       └── hello_world.pb.gw.go   # 生成的网关代码
├── route/                # HTTP路由
│   └── route.go
├── Makefile              # 构建自动化
├── Dockerfile            # 容器规范
├── go.mod                # 依赖项
└── README.md
```

## 🧪 测试

```bash
# 运行测试
make test

# 安装依赖
make deps
```

## 🤝 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启Pull Request

## 📄 许可证

MIT许可证 - 详情见[LICENSE](LICENSE)文件。