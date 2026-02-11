# TaskFlow

gRPC 任务调度服务。

## 🚀 特性

- **四种 gRPC 通信模式**
  - Simple RPC：单次请求/响应
  - Server Stream：服务端推送
  - Client Stream：批量创建任务
  - Bidirectional Stream：实时双向通信

- **异步任务处理**
  - SQLite 持久化存储
  - 任务队列管理
  - 状态机控制
  - 自动重试机制

- **生产级特性**
  - JWT 认证
  - 限流控制
  - 请求日志
  - 配置热加载

## 🏗️ 架构

```
taskflow/
├── cmd/
│   └── server/              # gRPC Server 入口
├── proto/
│   └── task.proto           # 服务定义
├── internal/
│   ├── config/              # 配置管理 (Worker, Queue, Database)
│   ├── handler/             # gRPC Handler
│   ├── middleware/          # 中间件 (认证, 限流, 日志)
│   ├── service/             # 业务逻辑
│   ├── repository/          # 数据访问层
│   ├── scheduler/           # 任务调度
│   ├── executor/            # 任务执行
│   └── queue/               # 消息队列
├── pkg/
│   └── errors/              # 错误定义
└── scripts/                 # 工具脚本
```

## 📦 技术栈

- **Go 1.21+**
- **gRPC** (Google Protocol Buffers)
- **SQLite** (持久化)
- **Zerolog** (日志)
- **JWT** (认证)

## 🛠️ 快速开始

```bash
# 安装依赖
make deps

# 构建项目
make build

# 运行服务
make run

# 运行测试
make test
```

## ⚙️ 配置

通过 `config.yaml` 或环境变量配置：

| 环境变量 | 描述 | 默认值 |
|---------|------|--------|
| GRPC_PORT | gRPC 端口 | 8080 |
| HTTP_PORT | HTTP 端口 | 8090 |
| DB_PATH | 数据库路径 | data/taskflow.db |
| WORKER_COUNT | Worker 数量 | 4 |
| MAX_RETRIES | 最大重试次数 | 3 |

## 📡 API 文档

### Simple RPC

```protobuf
service TaskService {
    rpc CreateTask(CreateTaskRequest) returns (Task);
    rpc GetTask(GetTaskRequest) returns (Task);
    rpc ListTasks(ListTasksRequest) returns (ListTasksResponse);
    rpc UpdateTask(UpdateTaskRequest) returns (Task);
    rpc DeleteTask(DeleteTaskRequest) returns (DeleteTaskResponse);
}
```

### Server Stream

```protobuf
rpc WatchTask(WatchTaskRequest) returns (stream TaskUpdate);
```

### Client Stream

```protobuf
rpc BatchCreateTasks(stream CreateTaskRequest) returns (BatchCreateResponse);
```

### Bidirectional Stream

```protobuf
rpc TaskUpdates(stream TaskCommand) returns (stream TaskUpdate);
```

## 📝 任务状态

| 状态 | 描述 |
|------|------|
| PENDING | 等待执行 |
| RUNNING | 执行中 |
| SUCCESS | 执行成功 |
| FAILED | 执行失败 |
| CANCELLED | 已取消 |

## 🧪 测试

```bash
# 单元测试
go test ./... -v

# 覆盖率
go test ./... -cover
```

## 📄 许可证

MIT

---

## 📌 更新日志

### v0.1.0 (2026-02-11)
- ✅ 项目初始化
- ✅ 配置系统扩展 (WorkerConfig, QueueConfig, DatabaseConfig)
- ✅ 完整配置验证逻辑
- ⏳ Task 数据模型 (待实现)
- ⏳ SQLite 持久化层 (待实现)
- ⏳ 错误处理 (待实现)
- ⏳ Handler 层 (待实现)
- ⏳ 流式 Handler (待实现)
- ⏳ 中间件 (待实现)
- ⏳ 服务层 (待实现)
- ⏳ 集成测试与文档 (待实现)
