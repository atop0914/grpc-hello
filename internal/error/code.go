package errorcode

// ErrorCode 错误码定义
type ErrorCode int32

const (
	// 通用错误 (1xxx)
	ErrCodeSuccess         ErrorCode = 0     // 成功
	ErrCodeUnknown         ErrorCode = 1000  // 未知错误
	ErrCodeInvalidParam    ErrorCode = 1001  // 参数错误
	ErrCodeUnauthorized    ErrorCode = 1002  // 未授权
	ErrCodeForbidden       ErrorCode = 1003 // 禁止访问
	ErrCodeNotFound        ErrorCode = 1004  // 资源不存在
	ErrCodeAlreadyExists   ErrorCode = 1005  // 资源已存在
	ErrCodeInvalidState    ErrorCode = 1006  // 状态无效
	ErrCodeTimeout         ErrorCode = 1007  // 超时
	ErrCodeRateLimit       ErrorCode = 1008  // 限流

	// 任务相关错误 (2xxx)
	ErrCodeTaskNotFound       ErrorCode = 2000 // 任务不存在
	ErrCodeTaskAlreadyRunning ErrorCode = 2001 // 任务已在运行
	ErrCodeTaskTerminated      ErrorCode = 2002 // 任务已终止
	ErrCodeTaskCancelled       ErrorCode = 2003 // 任务已取消
	ErrCodeTaskTimeout         ErrorCode = 2004 // 任务执行超时
	ErrCodeTaskDependency      ErrorCode = 2005 // 任务依赖未满足
	ErrCodeTaskRetryExhausted  ErrorCode = 2006 // 重试次数耗尽

	// 存储相关错误 (3xxx)
	ErrCodeDBError        ErrorCode = 3000 // 数据库错误
	ErrCodeDBNotConnected ErrorCode = 3001 // 数据库未连接
	ErrCodeDBTransaction  ErrorCode = 3002 // 事务错误

	// gRPC 相关错误 (4xxx)
	ErrCodeGRPCNotReady   ErrorCode = 4000 // gRPC 服务未就绪
	ErrCodeGRPCConnection ErrorCode = 4001 // gRPC 连接错误
	ErrCodeGRPCDeadline   ErrorCode = 4002 // gRPC 超时
)

// ErrorCodeMap 错误码到错误消息的映射
var ErrorCodeMap = map[ErrorCode]string{
	// 通用错误
	ErrCodeSuccess:        "success",
	ErrCodeUnknown:        "unknown error",
	ErrCodeInvalidParam:   "invalid parameter",
	ErrCodeUnauthorized:   "unauthorized",
	ErrCodeForbidden:      "forbidden",
	ErrCodeNotFound:       "resource not found",
	ErrCodeAlreadyExists:  "resource already exists",
	ErrCodeInvalidState:   "invalid state",
	ErrCodeTimeout:        "timeout",
	ErrCodeRateLimit:      "rate limit exceeded",

	// 任务相关
	ErrCodeTaskNotFound:       "task not found",
	ErrCodeTaskAlreadyRunning: "task already running",
	ErrCodeTaskTerminated:     "task already terminated",
	ErrCodeTaskCancelled:     "task cancelled",
	ErrCodeTaskTimeout:        "task timeout",
	ErrCodeTaskDependency:    "task dependency not satisfied",
	ErrCodeTaskRetryExhausted: "task retry exhausted",

	// 存储相关
	ErrCodeDBError:        "database error",
	ErrCodeDBNotConnected: "database not connected",
	ErrCodeDBTransaction:  "database transaction error",

	// gRPC 相关
	ErrCodeGRPCNotReady:   "gRPC service not ready",
	ErrCodeGRPCConnection: "gRPC connection error",
	ErrCodeGRPCDeadline:  "gRPC deadline exceeded",
}

// GetCodeMsg 获取错误码对应的消息
func GetCodeMsg(code ErrorCode) string {
	if msg, ok := ErrorCodeMap[code]; ok {
		return msg
	}
	return "unknown error"
}
