package errorcode

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TaskError 任务服务错误结构
type TaskError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Detail     string    `json:"detail,omitempty"`
	HTTPStatus int       `json:"-"`
}

// Error 实现 error 接口
func (e *TaskError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 解包错误
func (e *TaskError) Unwrap() error {
	return errors.New(e.Message)
}

// NewTaskError 创建新的任务错误
func NewTaskError(code ErrorCode, detail string) *TaskError {
	return &TaskError{
		Code:       code,
		Message:    GetCodeMsg(code),
		Detail:     detail,
		HTTPStatus: HTTPStatusFromCode(code),
	}
}

// NewTaskErrorWithMsg 创建带自定义消息的错误
func NewTaskErrorWithMsg(code ErrorCode, msg, detail string) *TaskError {
	return &TaskError{
		Code:       code,
		Message:    msg,
		Detail:     detail,
		HTTPStatus: HTTPStatusFromCode(code),
	}
}

// HTTPStatusFromCode 将错误码转换为 HTTP 状态码
func HTTPStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeSuccess:
		return http.StatusOK
	case ErrCodeInvalidParam:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeTaskNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists:
		return http.StatusConflict
	case ErrCodeInvalidState, ErrCodeTaskAlreadyRunning, ErrCodeTaskTerminated, ErrCodeTaskCancelled:
		return http.StatusBadRequest
	case ErrCodeTimeout, ErrCodeTaskTimeout, ErrCodeGRPCDeadline:
		return http.StatusGatewayTimeout
	case ErrCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrCodeDBError, ErrCodeDBNotConnected, ErrCodeDBTransaction, ErrCodeUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ToGRPCStatus 将 TaskError 转换为 gRPC status
func (e *TaskError) ToGRPCStatus() *status.Status {
	switch e.Code {
	case ErrCodeSuccess:
		return status.New(codes.OK, e.Message)
	case ErrCodeInvalidParam:
		return status.New(codes.InvalidArgument, e.Message)
	case ErrCodeUnauthorized:
		return status.New(codes.Unauthenticated, e.Message)
	case ErrCodeForbidden:
		return status.New(codes.PermissionDenied, e.Message)
	case ErrCodeNotFound, ErrCodeTaskNotFound:
		return status.New(codes.NotFound, e.Message)
	case ErrCodeAlreadyExists:
		return status.New(codes.AlreadyExists, e.Message)
	case ErrCodeTimeout, ErrCodeTaskTimeout:
		return status.New(codes.DeadlineExceeded, e.Message)
	case ErrCodeRateLimit:
		return status.New(codes.ResourceExhausted, e.Message)
	case ErrCodeDBError, ErrCodeDBNotConnected, ErrCodeDBTransaction:
		return status.New(codes.Internal, e.Message)
	case ErrCodeGRPCNotReady, ErrCodeGRPCConnection:
		return status.New(codes.Unavailable, e.Message)
	default:
		return status.New(codes.Unknown, e.Message)
	}
}

// FromGRPCStatus 从 gRPC status 创建 TaskError
func FromGRPCStatus(s *status.Status) *TaskError {
	code := ErrCodeUnknown
	httpStatus := http.StatusInternalServerError

	switch s.Code() {
	case codes.OK:
		code = ErrCodeSuccess
		httpStatus = http.StatusOK
	case codes.InvalidArgument:
		code = ErrCodeInvalidParam
		httpStatus = http.StatusBadRequest
	case codes.Unauthenticated:
		code = ErrCodeUnauthorized
		httpStatus = http.StatusUnauthorized
	case codes.PermissionDenied:
		code = ErrCodeForbidden
		httpStatus = http.StatusForbidden
	case codes.NotFound:
		code = ErrCodeNotFound
		httpStatus = http.StatusNotFound
	case codes.AlreadyExists:
		code = ErrCodeAlreadyExists
		httpStatus = http.StatusConflict
	case codes.DeadlineExceeded:
		code = ErrCodeTimeout
		httpStatus = http.StatusGatewayTimeout
	case codes.ResourceExhausted:
		code = ErrCodeRateLimit
		httpStatus = http.StatusTooManyRequests
	case codes.Internal:
		code = ErrCodeDBError
		httpStatus = http.StatusInternalServerError
	case codes.Unavailable:
		code = ErrCodeGRPCNotReady
		httpStatus = http.StatusServiceUnavailable
	}

	return &TaskError{
		Code:       code,
		Message:    s.Message(),
		HTTPStatus: httpStatus,
	}
}

// GinErrorResponse Gin 错误响应结构
type GinErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

// ToGinResponse 转换为 Gin JSON 响应
func (e *TaskError) ToGinResponse() GinErrorResponse {
	return GinErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Detail:  e.Detail,
	}
}

// HandleGinError 处理 Gin 错误响应
func HandleGinError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var taskErr *TaskError
	if errors.As(err, &taskErr) {
		c.JSON(taskErr.HTTPStatus, taskErr.ToGinResponse())
		return
	}

	// 未知错误
	taskErr = NewTaskError(ErrCodeUnknown, err.Error())
	c.JSON(taskErr.HTTPStatus, taskErr.ToGinResponse())
}

// HandleGinErrorWithCode 使用指定错误码处理错误
func HandleGinErrorWithCode(c *gin.Context, code ErrorCode, detail string) {
	err := NewTaskError(code, detail)
	HandleGinError(c, err)
}

// HandleGinPanic 处理 Panic
func HandleGinPanic(c *gin.Context, recovered interface{}) {
	var detail string
	if recovered != nil {
		detail = fmt.Sprintf("%v", recovered)
	}

	err := NewTaskError(ErrCodeUnknown, "internal server error")
	if detail != "" {
		err.Detail = detail
	}

	c.JSON(err.HTTPStatus, err.ToGinResponse())
}
