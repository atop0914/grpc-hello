package handler

import "log"

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	errorcode "taskflow/internal/error"
	"taskflow/internal/model"
	"taskflow/internal/repository"
	pb "taskflow/proto"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	repo         *repository.TaskRepository
	watchers     map[string][]chan *pb.TaskChangeEvent
	watchersMu   sync.RWMutex
	taskUpdateCh chan *pb.TaskChangeEvent
	pb.UnimplementedTaskServiceServer
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	h := &TaskHandler{
		repo:         repo,
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 100),
	}
	// 启动任务变更通知循环
	go h.taskUpdateNotifier()
	return h
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	// 参数验证
	if req.Name == "" {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeInvalidParam, "name is required").ToGRPCStatus().Err()
	}

	// 创建任务模型
	task := model.NewTask(
		req.Name,
		req.Description,
		model.TaskPriority(req.Priority),
		req.TaskType,
		req.InputParams,
		req.Dependencies,
		req.MaxRetries,
		req.CreatedBy,
	)
	task.ID = uuid.New().String()

	// 保存到数据库
	if err := h.repo.Create(task); err != nil {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeDBError, err.Error()).ToGRPCStatus().Err()
	}

	return h.toPBTask(task, false), nil
}

// GetTask 获取任务
func (h *TaskHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	if req.Id == "" {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeInvalidParam, "id is required").ToGRPCStatus().Err()
	}

	task, err := h.repo.GetByID(req.Id)
	if err != nil { log.Printf("Handler error: %v", err)
		return nil, errorcode.NewTaskError(errorcode.ErrCodeDBError, err.Error()).ToGRPCStatus().Err()
	}
	if task == nil {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeTaskNotFound, "task not found").ToGRPCStatus().Err()
	}

	return h.toPBTask(task, req.IncludeEvents), nil
}

// ListTasks 列出任务
func (h *TaskHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	// 分页参数
	pageSize := int(req.PageSize)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	pageIndex := int(req.Page)
	if pageIndex < 0 {
		pageIndex = 0
	}
	offset := pageIndex * pageSize

	// 构建过滤条件
	filter := repository.TaskFilter{
		PageSize:  pageSize,
		PageIndex: offset,
		Keyword:   req.Keyword,
		TaskType:  req.TaskType,
	}

	if len(req.StatusFilter) > 0 {
		status := model.TaskStatus(req.StatusFilter[0])
		filter.Status = &status
	}
	if req.Priority != 0 {
		priority := model.TaskPriority(req.Priority)
		filter.Priority = &priority
	}

	// 查询
	tasks, total, err := h.repo.ListByFilter(filter)
	if err != nil { log.Printf("Handler error: %v", err)
		return nil, errorcode.NewTaskError(errorcode.ErrCodeDBError, err.Error()).ToGRPCStatus().Err()
	}

	// 转换
	pbTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		pbTasks[i] = h.toPBTask(task, false)
	}

	return &pb.ListTasksResponse{
		Tasks:    pbTasks,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.Task, error) {
	if req.Id == "" {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeInvalidParam, "id is required").ToGRPCStatus().Err()
	}

	// 获取现有任务
	task, err := h.repo.GetByID(req.Id)
	if err != nil { log.Printf("Handler error: %v", err)
		return nil, errorcode.NewTaskError(errorcode.ErrCodeDBError, err.Error()).ToGRPCStatus().Err()
	}
	if task == nil {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeTaskNotFound, "task not found").ToGRPCStatus().Err()
	}

	// 更新字段
	if req.Status != 0 {
		oldStatus := task.Status
		newStatus := model.TaskStatus(req.Status)

		// 状态转换验证
		if !isValidStatusTransition(oldStatus, newStatus) {
			return nil, errorcode.NewTaskError(errorcode.ErrCodeInvalidState,
				fmt.Sprintf("invalid status transition from %s to %s", oldStatus, newStatus)).ToGRPCStatus().Err()
		}

		// 原子更新状态
		err := h.repo.UpdateStatusWithEvent(req.Id, oldStatus, newStatus, "system", "status updated")
		if err != nil { log.Printf("Handler error: %v", err)
			return nil, errorcode.NewTaskError(errorcode.ErrCodeDBError, err.Error()).ToGRPCStatus().Err()
		}
		task.Status = newStatus
	}

	if req.OutputResult != nil {
		task.OutputResult = req.OutputResult
	}
	if req.ErrorMessage != "" {
		task.ErrorMessage = req.ErrorMessage
	}
	if req.RetryCount != 0 {
		task.RetryCount = req.RetryCount
	}
	task.UpdatedAt = time.Now()

	// 保存
	if err := h.repo.Update(task); err != nil {
		return nil, errorcode.NewTaskError(errorcode.ErrCodeDBError, err.Error()).ToGRPCStatus().Err()
	}

	return h.toPBTask(task, false), nil
}

// 状态转换验证
func isValidStatusTransition(from, to model.TaskStatus) bool {
	// PENDING 可以转到 RUNNING, CANCELLED
	if from == model.TaskStatusPending {
		return to == model.TaskStatusRunning || to == model.TaskStatusCancelled
	}
	// RUNNING 可以转到 SUCCEEDED, FAILED, TIMEOUT, CANCELLED
	if from == model.TaskStatusRunning {
		return to == model.TaskStatusSucceeded ||
			to == model.TaskStatusFailed ||
			to == model.TaskStatusTimeout ||
			to == model.TaskStatusCancelled
	}
	// 终态不能转换
	return false
}

// toPBTask 转换为 Protobuf 任务
func (h *TaskHandler) toPBTask(task *model.Task, includeEvents bool) *pb.Task {
	pbTask := &pb.Task{
		Id:           task.ID,
		Name:         task.Name,
		Description:  task.Description,
		Status:       pb.TaskStatus(task.Status),
		Priority:     pb.TaskPriority(task.Priority),
		TaskType:     task.TaskType,
		InputParams:  task.InputParams,
		OutputResult: task.OutputResult,
		Dependencies: task.Dependencies,
		RetryCount:   task.RetryCount,
		MaxRetries:   task.MaxRetries,
		ErrorMessage: task.ErrorMessage,
		CreatedAt:    task.CreatedAt.Unix(),
		UpdatedAt:    task.UpdatedAt.Unix(),
		CreatedBy:    task.CreatedBy,
	}

	if task.StartedAt != nil {
		pbTask.StartedAt = task.StartedAt.Unix()
	}
	if task.CompletedAt != nil {
		pbTask.CompletedAt = task.CompletedAt.Unix()
	}

	if includeEvents {
		for _, e := range task.Events {
			pbTask.Events = append(pbTask.Events, &pb.TaskEvent{
				Id:         e.ID,
				FromStatus: pb.TaskStatus(e.FromStatus),
				ToStatus:   pb.TaskStatus(e.ToStatus),
				Message:    e.Message,
				Timestamp:  e.Timestamp.Unix(),
				Operator:   e.Operator,
			})
		}
	}

	return pbTask
}

// RegisterTaskHandlers 注册任务服务句柄
func RegisterTaskHandlers(repo *repository.TaskRepository) *TaskHandler {
	return NewTaskHandler(repo)
}

// ========== 流式 RPC 实现 ==========

// taskUpdateNotifier 任务变更通知器
func (h *TaskHandler) taskUpdateNotifier() {
	for event := range h.taskUpdateCh {
		h.notifyWatchers(event)
	}
}

// notifyWatchers 通知所有订阅者
func (h *TaskHandler) notifyWatchers(event *pb.TaskChangeEvent) {
	h.watchersMu.RLock()
	defer h.watchersMu.RUnlock()

	if chs, ok := h.watchers[event.TaskId]; ok {
		for _, ch := range chs {
			select {
			case ch <- event:
			default:
			}
		}
	}

	if globalChs, ok := h.watchers[""]; ok {
		for _, ch := range globalChs {
			select {
			case ch <- event:
			default:
			}
		}
	}
}

// broadcastTaskChange 广播任务变更
func (h *TaskHandler) broadcastTaskChange(taskId string, task *model.Task, fromStatus, toStatus model.TaskStatus, changeType string) {
	event := &pb.TaskChangeEvent{
		TaskId:     taskId,
		Task:       h.toPBTask(task, false),
		FromStatus: pb.TaskStatus(fromStatus),
		ToStatus:   pb.TaskStatus(toStatus),
		ChangedAt:  time.Now().Unix(),
		ChangeType: changeType,
	}
	h.taskUpdateCh <- event
}

// WatchTask 服务端流式 - 监听任务状态变化
func (h *TaskHandler) WatchTask(req *pb.WatchTaskRequest, stream pb.TaskService_WatchTaskServer) error {
	ch := make(chan *pb.TaskChangeEvent, 10)
	taskIDs := req.TaskIds

	h.watchersMu.Lock()
	watchKey := ""
	if len(taskIDs) == 1 {
		watchKey = taskIDs[0]
	}
	h.watchers[watchKey] = append(h.watchers[watchKey], ch)
	h.watchersMu.Unlock()

	if req.IncludeInitial {
		var tasks []*model.Task
		if len(taskIDs) > 0 {
			for _, id := range taskIDs {
				task, err := h.repo.GetByID(id)
				if err == nil && task != nil {
					tasks = append(tasks, task)
				}
			}
		} else {
			tasks, _, _ = h.repo.ListByFilter(repository.TaskFilter{PageSize: 50, PageIndex: 0})
		}

		for _, task := range tasks {
			event := &pb.TaskChangeEvent{
				TaskId:     task.ID,
				Task:       h.toPBTask(task, false),
				FromStatus: pb.TaskStatus(task.Status),
				ToStatus:   pb.TaskStatus(task.Status),
				ChangedAt:  task.UpdatedAt.Unix(),
				ChangeType: "initial",
			}
			stream.Send(event)
		}
	}

	ctx := stream.Context()
	defer func() {
		h.watchersMu.Lock()
		if chs, ok := h.watchers[watchKey]; ok {
			for i, c := range chs {
				if c == ch {
					h.watchers[watchKey] = append(chs[:i], chs[i+1:]...)
					break
				}
			}
		}
		h.watchersMu.Unlock()
		close(ch)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-ch:
			if len(req.StatusFilter) > 0 {
				filtered := true
				for _, s := range req.StatusFilter {
					if event.ToStatus == s {
						filtered = false
						break
					}
				}
				if filtered {
					continue
				}
			}
			stream.Send(event)
		}
	}
}

// BatchCreateTasks 客户端流式 - 批量创建任务
func (h *TaskHandler) BatchCreateTasks(stream pb.TaskService_BatchCreateTasksServer) error {
	var tasks []*pb.Task
	var errors []string
	successCount := 0
	failedCount := 0

	for {
		req, err := stream.Recv()
		if err != nil {
			break
		}

		if req.Name == "" {
			failedCount++
			errors = append(errors, "name is required")
			tasks = append(tasks, nil)
			continue
		}

		task := model.NewTask(
			req.Name,
			req.Description,
			model.TaskPriority(req.Priority),
			req.TaskType,
			req.InputParams,
			req.Dependencies,
			req.MaxRetries,
			req.CreatedBy,
		)
		task.ID = uuid.New().String()

		if err := h.repo.Create(task); err != nil {
			failedCount++
			errors = append(errors, err.Error())
			tasks = append(tasks, nil)
			continue
		}

		h.broadcastTaskChange(task.ID, task, model.TaskStatusUnspecified, model.TaskStatusPending, "created")

		pbTask := h.toPBTask(task, false)
		tasks = append(tasks, pbTask)
		successCount++
	}

	return stream.SendAndClose(&pb.BatchCreateTasksResponse{
		Tasks:        tasks,
		SuccessCount: int32(successCount),
		FailedCount:  int32(failedCount),
		Errors:       errors,
	})
}

// TaskUpdates 双向流式 - 任务更新流
func (h *TaskHandler) TaskUpdates(stream pb.TaskService_TaskUpdatesServer) error {
	ctx := stream.Context()
	wg := sync.WaitGroup{}

	eventCh := make(chan *pb.TaskChangeEvent, 10)
	recvCh := make(chan *pb.TaskUpdateRequest, 10)
	sendCh := make(chan *pb.TaskUpdateResponse, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err != nil {
				close(recvCh)
				return
			}
			recvCh <- req
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for resp := range sendCh {
			stream.Send(resp)
		}
	}()

	globalCh := make(chan *pb.TaskChangeEvent, 10)
	h.watchersMu.Lock()
	h.watchers[""] = append(h.watchers[""], globalCh)
	h.watchersMu.Unlock()

	defer func() {
		h.watchersMu.Lock()
		if chs, ok := h.watchers[""]; ok {
			for i, c := range chs {
				if c == globalCh {
					h.watchers[""] = append(chs[:i], chs[i+1:]...)
					break
				}
			}
		}
		h.watchersMu.Unlock()
		close(eventCh)
		close(sendCh)
		wg.Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req := <-recvCh:
			requestId := req.GetRequestId()
			updateType := req.GetUpdateType()

			switch updateType {
			case "create":
				createReq := req.GetCreate()
				if createReq != nil {
					task, err := h.CreateTask(ctx, createReq)
					resp := &pb.TaskUpdateResponse{
						RequestId: requestId,
						Success:   err == nil,
						Task:      task,
					}
					if err != nil {
						resp.Error = err.Error()
					}
					sendCh <- resp
				}
			case "update":
				updateReq := req.GetUpdate()
				if updateReq != nil {
					task, err := h.UpdateTask(ctx, updateReq)
					resp := &pb.TaskUpdateResponse{
						RequestId: requestId,
						Success:   err == nil,
						Task:      task,
					}
					if err != nil {
						resp.Error = err.Error()
					}
					sendCh <- resp
				}
			default:
				sendCh <- &pb.TaskUpdateResponse{
					RequestId: requestId,
					Success:   false,
					Error:     "unknown update type",
				}
			}
		case event := <-globalCh:
			resp := &pb.TaskUpdateResponse{
				ChangeEvent: event,
				Success:     true,
			}
			sendCh <- resp
		}
	}
}
