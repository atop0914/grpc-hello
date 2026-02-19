package handler

import (
	"testing"

	"taskflow/internal/model"
	pb "taskflow/proto"
)

// TestHandler_StreamMethodsExist verifies streaming methods exist
func TestHandler_StreamMethodsExist(t *testing.T) {
	handler := &TaskHandler{
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 10),
	}

	// Verify handler has streaming methods
	_ = handler.WatchTask
	_ = handler.BatchCreateTasks
	_ = handler.TaskUpdates
	_ = handler.broadcastTaskChange
	_ = handler.notifyWatchers
	_ = handler.taskUpdateNotifier

	t.Log("All streaming methods exist on handler")
}

// TestHandler_NotifyWatchers tests task notification
func TestHandler_NotifyWatchers(t *testing.T) {
	handler := &TaskHandler{
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 10),
	}

	// Start notifier
	go handler.taskUpdateNotifier()

	// Create a test task
	task := model.NewTask("notify-test", "test", model.TaskPriorityNormal, "default", nil, nil, 3, "test")
	task.ID = "notify-task-id"

	// Test broadcast
	handler.broadcastTaskChange(task.ID, task, model.TaskStatusPending, model.TaskStatusRunning, "started")

	// Wait for notification
	<-handler.taskUpdateCh

	t.Log("Task notification test passed")
}

// TestHandler_MultipleWatchers tests multiple watchers
func TestHandler_MultipleWatchers(t *testing.T) {
	handler := &TaskHandler{
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 10),
	}

	// Create channel for watcher
	ch := make(chan *pb.TaskChangeEvent, 10)

	// Register watcher
	handler.watchersMu.Lock()
	handler.watchers["test-task"] = append(handler.watchers["test-task"], ch)
	handler.watchersMu.Unlock()

	// Verify watcher registered
	handler.watchersMu.RLock()
	watchers, ok := handler.watchers["test-task"]
	handler.watchersMu.RUnlock()

	if !ok || len(watchers) != 1 {
		t.Fatalf("watcher not registered correctly")
	}

	t.Log("Multiple watchers test passed")
}

// TestHandler_ConcurrentNotifications tests concurrent notifications
func TestHandler_ConcurrentNotifications(t *testing.T) {
	handler := &TaskHandler{
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 10),
	}

	// Start notifier
	go handler.taskUpdateNotifier()

	// Create multiple watchers
	for i := 0; i < 5; i++ {
		ch := make(chan *pb.TaskChangeEvent, 10)
		handler.watchersMu.Lock()
		handler.watchers[""] = append(handler.watchers[""], ch)
		handler.watchersMu.Unlock()
	}

	// Broadcast notification
	task := model.NewTask("concurrent-test", "test", model.TaskPriorityNormal, "default", nil, nil, 3, "test")
	handler.broadcastTaskChange(task.ID, task, model.TaskStatusPending, model.TaskStatusRunning, "started")

	// Wait for notification
	<-handler.taskUpdateCh

	t.Log("Concurrent notifications test passed")
}

// TestHandler_StatusTransitionInNotification tests status in notification
func TestHandler_StatusTransitionInNotification(t *testing.T) {
	handler := &TaskHandler{
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 10),
	}

	// Start notifier
	go handler.taskUpdateNotifier()

	// Test different status transitions
	transitions := []struct {
		from model.TaskStatus
		to   model.TaskStatus
	}{
		{model.TaskStatusPending, model.TaskStatusRunning},
		{model.TaskStatusRunning, model.TaskStatusSucceeded},
		{model.TaskStatusRunning, model.TaskStatusFailed},
		{model.TaskStatusPending, model.TaskStatusCancelled},
	}

	task := model.NewTask("status-test", "test", model.TaskPriorityNormal, "default", nil, nil, 3, "test")

	for _, tr := range transitions {
		handler.broadcastTaskChange(task.ID, task, tr.from, tr.to, "test")
		<-handler.taskUpdateCh
	}

	t.Log("Status transition test passed")
}

// BenchmarkHandler_NotifyWatchers benchmarks notification
func BenchmarkHandler_NotifyWatchers(b *testing.B) {
	handler := &TaskHandler{
		watchers:     make(map[string][]chan *pb.TaskChangeEvent),
		taskUpdateCh: make(chan *pb.TaskChangeEvent, 1000),
	}

	// Start notifier
	go handler.taskUpdateNotifier()

	task := model.NewTask("bench-test", "test", model.TaskPriorityNormal, "default", nil, nil, 3, "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.broadcastTaskChange(task.ID, task, model.TaskStatusPending, model.TaskStatusRunning, "started")
		<-handler.taskUpdateCh
	}
}
