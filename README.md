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
├── proto/
│   └── task.proto           # 服务定义
├── internal/
│   ├── config/             # 配置管理
│   ├── model/               # 数据模型
│   ├── error/               # 错误码定义与处理
│   ├── repository/          # SQLite 数据访问层 ✅ 已完成
│   ├── middleware/          # 中间件 ⏳ 待实现
│   ├── handler/            # gRPC Handler ⏳ 待实现
│   ├── service/            # 业务逻辑层 ⏳ 待实现
│   ├── scheduler/           # 任务调度 ⏳ 待实现
│   ├── executor/           # 任务执行 ⏳ 待实现
│   └── queue/              # 消息队列 ⏳ 待实现
├── cmd/
│   └── server/              # 服务入口
└── scripts/                 # 工具脚本
```

## 📦 技术栈

- **Go 1.21+**
- **gRPC** (Google Protocol Buffers)
- **SQLite** (持久化)
- **Zerolog** (日志)

## 🛠️ 快速开始

```bash
# 安装依赖
go mod tidy

# 构建项目
go build -o taskflow .

# 运行服务
./taskflow

# 运行测试
go test ./...
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

## ✅ 已完成功能

### 1. SQLite 持久化层 (internal/repository/)

提供完整的 CRUD 操作：

| 方法 | 描述 |
|------|------|
| `Create` | 创建任务 |
| `GetByID` | 根据 ID 获取任务 |
| `Update` | 更新任务 |
| `Delete` | 删除任务 |
| `List` | 分页列出任务 |
| `ListByStatus` | 按状态列出任务 |
| `ListPending` | 列出待处理任务 |
| `ListByCreator` | 按创建者查询 |
| `ListByFilter` | 多条件过滤查询 |
| `Search` | 关键词搜索 |
| `Count` | 统计任务数量 |
| `UpdateStatus` | 更新任务状态 |
| `UpdateStatusWithEvent` | 原子更新+记录事件 |
| `AddEvent` | 添加任务事件 |
| `GetEventsByTaskID` | 获取任务所有事件 |

### 2. 错误处理模块 (internal/error/)

完整的错误码定义和错误处理函数：

**错误码定义：**
- 通用错误 (1xxx)：参数错误、未授权、禁止访问、未找到、超时等
- 任务相关错误 (2xxx)：任务未找到、运行中、终止/取消/超时、依赖未满足等
- 存储相关错误 (3xxx)：数据库错误、未连接、事务错误
- gRPC 相关错误 (4xxx)：服务未就绪、连接错误、超时

**错误处理函数：**
- `TaskError` 结构体实现 error 接口
- `HTTPStatusFromCode()` - 错误码转 HTTP 状态码
- `ToGRPCStatus()` / `FromGRPCStatus()` - gRPC status 互转
- `HandleGinError()` / `HandleGinErrorWithCode()` - Gin 中间件错误处理
- `HandleGinPanic()` - Panic 恢复处理

## ⏳ 待实现功能

- [ ] 中间件 (认证, 限流, 日志)
- [ ] Handler 层 (Simple RPC)
- [ ] 流式 Handler (Server/Client/Bidirectional Stream)
- [ ] 服务层 (任务调度, 状态机)
- [ ] 集成测试与文档

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

### v0.1.0 (2026-02-14)
- ✅ 项目初始化
- ✅ 配置系统扩展 (WorkerConfig, QueueConfig, DatabaseConfig)
- ✅ 完整配置验证逻辑
- ✅ Task 数据模型
- ✅ SQLite 持久化层 (Repository)
- ✅ 错误处理模块
- ⏳ Handler 层
- ⏳ 流式 Handler
- ⏳ 中间件
- ⏳ 服务层
- ⏳ 集成测试与文档
