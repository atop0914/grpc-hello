package model

import (
	"time"
)

// TaskStatus 任务状态枚举
type TaskStatus int32

const (
	TaskStatusUnspecified TaskStatus = 0
	TaskStatusPending     TaskStatus = 1
	TaskStatusRunning     TaskStatus = 2
	TaskStatusSucceeded   TaskStatus = 3
	TaskStatusFailed      TaskStatus = 4
	TaskStatusCancelled   TaskStatus = 5
	TaskStatusTimeout     TaskStatus = 6
)

func (s TaskStatus) String() string {
	switch s {
	case TaskStatusPending:
		return "PENDING"
	case TaskStatusRunning:
		return "RUNNING"
	case TaskStatusSucceeded:
		return "SUCCEEDED"
	case TaskStatusFailed:
		return "FAILED"
	case TaskStatusCancelled:
		return "CANCELLED"
	case TaskStatusTimeout:
		return "TIMEOUT"
	default:
		return "UNSPECIFIED"
	}
}

// TaskPriority 任务优先级枚举
type TaskPriority int32

const (
	TaskPriorityUnspecified TaskPriority = 0
	TaskPriorityLow         TaskPriority = 1
	TaskPriorityNormal      TaskPriority = 2
	TaskPriorityHigh        TaskPriority = 3
	TaskPriorityUrgent      TaskPriority = 4
)

func (p TaskPriority) String() string {
	switch p {
	case TaskPriorityLow:
		return "LOW"
	case TaskPriorityNormal:
		return "NORMAL"
	case TaskPriorityHigh:
		return "HIGH"
	case TaskPriorityUrgent:
		return "URGENT"
	default:
		return "UNSPECIFIED"
	}
}

// Task 任务实体
type Task struct {
	ID            string            `json:"id" bson:"_id"`
	Name          string            `json:"name" bson:"name"`
	Description   string            `json:"description" bson:"description"`
	Status        TaskStatus        `json:"status" bson:"status"`
	Priority      TaskPriority      `json:"priority" bson:"priority"`
	TaskType      string            `json:"task_type" bson:"task_type"`
	InputParams   map[string]string `json:"input_params" bson:"input_params"`
	OutputResult  map[string]string `json:"output_result" bson:"output_result"`
	Dependencies  []string          `json:"dependencies" bson:"dependencies"`
	RetryCount    int32             `json:"retry_count" bson:"retry_count"`
	MaxRetries    int32             `json:"max_retries" bson:"max_retries"`
	ErrorMessage  string            `json:"error_message" bson:"error_message"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" bson:"updated_at"`
	StartedAt     *time.Time        `json:"started_at,omitempty" bson:"started_at,omitempty"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	CreatedBy     string            `json:"created_by" bson:"created_by"`
	Events        []TaskEvent       `json:"events" bson:"events"`
}

// TaskEvent 任务状态变更事件
type TaskEvent struct {
	ID         string     `json:"id" bson:"_id"`
	TaskID     string     `json:"task_id" bson:"task_id"`
	FromStatus TaskStatus `json:"from_status" bson:"from_status"`
	ToStatus   TaskStatus `json:"to_status" bson:"to_status"`
	Message    string     `json:"message" bson:"message"`
	Timestamp  time.Time  `json:"timestamp" bson:"timestamp"`
	Operator   string     `json:"operator" bson:"operator"`
}

// NewTask 创建新任务
func NewTask(name, description string, priority TaskPriority, taskType string, inputParams map[string]string, dependencies []string, maxRetries int32, createdBy string) *Task {
	now := time.Now()
	return &Task{
		Name:         name,
		Description:  description,
		Priority:     priority,
		TaskType:     taskType,
		InputParams:  inputParams,
		Dependencies: dependencies,
		MaxRetries:   maxRetries,
		CreatedBy:    createdBy,
		Status:       TaskStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsTerminal 检查任务是否处于终态
func (t *Task) IsTerminal() bool {
	return t.Status == TaskStatusSucceeded ||
		t.Status == TaskStatusFailed ||
		t.Status == TaskStatusCancelled ||
		t.Status == TaskStatusTimeout
}

// CanRetry 检查任务是否可重试
func (t *Task) CanRetry() bool {
	return t.Status == TaskStatusFailed && t.RetryCount < t.MaxRetries
}

// MarkRunning 标记任务为运行中
func (t *Task) MarkRunning() {
	t.Status = TaskStatusRunning
	now := time.Now()
	t.StartedAt = &now
	t.UpdatedAt = now
}

// MarkCompleted 标记任务为完成
func (t *Task) MarkCompleted() {
	t.Status = TaskStatusSucceeded
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

// MarkFailed 标记任务失败
func (t *Task) MarkFailed(errMsg string) {
	t.Status = TaskStatusFailed
	t.ErrorMessage = errMsg
	t.RetryCount++
	t.UpdatedAt = time.Now()
}
