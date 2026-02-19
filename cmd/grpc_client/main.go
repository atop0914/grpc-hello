package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "taskflow/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	addr = "localhost:8080"
)

func main() {
	fmt.Println("========== TaskFlow gRPC 客户端测试 ==========")
	
	// 连接 gRPC 服务
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()
	
	client := pb.NewTaskServiceClient(conn)
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// ========== 测试 1: Simple RPC ==========
	fmt.Println("\n【1】Simple RPC 测试")
	testSimpleRPC(ctx, client)
	
	// ========== 测试 2: Server Streaming ==========
	fmt.Println("\n【2】Server Streaming 测试")
	testServerStreaming(ctx, client)
	
	// ========== 测试 3: Client Streaming ==========
	fmt.Println("\n【3】Client Streaming 测试")
	testClientStreaming(ctx, client)
	
	// ========== 测试 4: Bidirectional Streaming ==========
	fmt.Println("\n【4】Bidirectional Streaming 测试")
	testBidirectionalStreaming(ctx, client)
	
	fmt.Println("\n========== 全部测试完成 ==========")
}

// ========== 1. Simple RPC 测试 ==========
func testSimpleRPC(ctx context.Context, client pb.TaskServiceClient) {
	// 1.1 CreateTask
	fmt.Println("  [1.1] CreateTask - 正常")
	task, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Name:        "grpc-test-1",
		Description: "Simple RPC 测试",
		Priority:    pb.TaskPriority(2), // NORMAL
	})
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	fmt.Printf("  ✅ 成功: ID=%s, Name=%s\n", task.Id, task.Name)
	
	// 1.2 CreateTask - 无name (应失败)
	fmt.Println("  [1.2] CreateTask - 无name (应失败)")
	_, err = client.CreateTask(ctx, &pb.CreateTaskRequest{
		Description: "无名称",
	})
	if err != nil {
		fmt.Printf("  ✅ 正确拒绝: %v\n", err)
	} else {
		fmt.Println("  ❌ 应该失败但没有")
	}
	
	// 1.3 GetTask
	fmt.Println("  [1.3] GetTask - 正常")
	task, err = client.GetTask(ctx, &pb.GetTaskRequest{
		Id: task.Id,
	})
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	fmt.Printf("  ✅ 成功: Status=%s\n", task.Status)
	
	// 1.4 GetTask - 不存在
	fmt.Println("  [1.4] GetTask - 不存在")
	_, err = client.GetTask(ctx, &pb.GetTaskRequest{
		Id: "not-exist",
	})
	if err != nil {
		fmt.Printf("  ✅ 正确返回错误: %v\n", err)
	} else {
		fmt.Println("  ❌ 应该失败但没有")
	}
	
	// 1.5 ListTasks
	fmt.Println("  [1.5] ListTasks")
	resp, err := client.ListTasks(ctx, &pb.ListTasksRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	fmt.Printf("  ✅ 成功: Total=%d\n", resp.Total)
	
	// 1.6 UpdateTask - PENDING->RUNNING
	fmt.Println("  [1.6] UpdateTask - PENDING->RUNNING")
	task, err = client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:     task.Id,
		Status: pb.TaskStatus(2), // RUNNING
	})
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	fmt.Printf("  ✅ 成功: Status=%s\n", task.Status)
	
	// 1.7 UpdateTask - RUNNING->SUCCEEDED
	fmt.Println("  [1.7] UpdateTask - RUNNING->SUCCEEDED")
	task, err = client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:          task.Id,
		Status:      pb.TaskStatus(3), // SUCCEEDED
		OutputResult: map[string]string{"result": "success"},
	})
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	fmt.Printf("  ✅ 成功: Status=%s\n", task.Status)
	
	// 1.8 UpdateTask - 非法转换 (SUCCEEDED->RUNNING)
	fmt.Println("  [1.8] UpdateTask - 非法转换")
	_, err = client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:     task.Id,
		Status: pb.TaskStatus(2), // RUNNING
	})
	if err != nil {
		fmt.Printf("  ✅ 正确拒绝: %v\n", err)
	} else {
		fmt.Println("  ❌ 应该失败但没有")
	}
	
	// 保存任务ID供后续测试使用
	testTaskId = task.Id
}

var testTaskId string

// ========== 2. Server Streaming 测试 ==========
func testServerStreaming(ctx context.Context, client pb.TaskServiceClient) {
	// 创建测试任务
	task1, _ := client.CreateTask(ctx, &pb.CreateTaskRequest{Name: "stream-test-1"})
	task2, _ := client.CreateTask(ctx, &pb.CreateTaskRequest{Name: "stream-test-2"})
	client.UpdateTask(ctx, &pb.UpdateTaskRequest{Id: task1.Id, Status: pb.TaskStatus(2)})
	
	// 2.1 WatchTask - 监听特定任务
	fmt.Println("  [2.1] WatchTask - 监听特定任务")
	stream, err := client.WatchTask(ctx, &pb.WatchTaskRequest{
		TaskIds:        []string{task1.Id, task2.Id},
		IncludeInitial: true,
	})
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	
	// 接收初始事件
	eventCount := 0
	for {
		event, err := stream.Recv()
		if err != nil {
			break
		}
		eventCount++
		fmt.Printf("  ✅ 收到事件 #%d: TaskId=%s, ChangeType=%s\n", 
			eventCount, event.TaskId, event.ChangeType)
		if eventCount >= 2 {
			break
		}
	}
	
	// 2.2 WatchTask - 全局监听
	fmt.Println("  [2.2] WatchTask - 全局监听")
	stream2, _ := client.WatchTask(ctx, &pb.WatchTaskRequest{
		IncludeInitial: false,
	})
	
	// 创建新任务触发事件
	client.CreateTask(ctx, &pb.CreateTaskRequest{Name: "trigger-task"})
	time.Sleep(200 * time.Millisecond)
	
	// 尝试接收
	select {
	case event, ok := <-recvWithTimeout(stream2):
		if ok {
			fmt.Printf("  ✅ 收到全局事件: %s\n", event.ChangeType)
		}
	default:
		fmt.Println("  ℹ️  无阻塞事件（正常）")
	}
	
	stream2.CloseSend()
	
	fmt.Println("  ✅ Server Streaming 测试完成")
}

func recvWithTimeout(stream pb.TaskService_WatchTaskClient) <-chan *pb.TaskChangeEvent {
	ch := make(chan *pb.TaskChangeEvent, 1)
	go func() {
		event, _ := stream.Recv()
		ch <- event
	}()
	return ch
}

// ========== 3. Client Streaming 测试 ==========
func testClientStreaming(ctx context.Context, client pb.TaskServiceClient) {
	// 3.1 BatchCreateTasks
	fmt.Println("  [3.1] BatchCreateTasks - 批量创建")
	stream, err := client.BatchCreateTasks(ctx)
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	
	// 发送多个请求
	tasks := []struct {
		Name        string
		Description string
		Priority    pb.TaskPriority
	}{
		{"batch-1", "批量任务1", pb.TaskPriority(2)},
		{"batch-2", "批量任务2", pb.TaskPriority(3)},
		{"batch-3", "批量任务3", pb.TaskPriority(1)},
	}
	
	for _, t := range tasks {
		req := &pb.CreateTaskRequest{
			Name:        t.Name,
			Description: t.Description,
			Priority:    t.Priority,
		}
		if err := stream.Send(req); err != nil {
			fmt.Printf("  ❌ 发送失败: %v\n", err)
			return
		}
		fmt.Printf("  ✅ 发送: %s\n", t.Name)
	}
	
	// 关闭并接收响应
	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("  ❌ 接收失败: %v\n", err)
		return
	}
	
	fmt.Printf("  ✅ 批量创建完成: Success=%d, Failed=%d\n", 
		resp.SuccessCount, resp.FailedCount)
	
	// 3.2 BatchCreateTasks - 包含错误请求
	fmt.Println("  [3.2] BatchCreateTasks - 混合错误")
	stream2, _ := client.BatchCreateTasks(ctx)
	stream2.Send(&pb.CreateTaskRequest{Name: "valid-1"})
	stream2.Send(&pb.CreateTaskRequest{}) // 无name - 应失败
	stream2.Send(&pb.CreateTaskRequest{Name: "valid-2"})
	resp2, _ := stream2.CloseAndRecv()
	fmt.Printf("  ✅ 结果: Success=%d, Failed=%d\n", 
		resp2.SuccessCount, resp2.FailedCount)
}

// ========== 4. Bidirectional Streaming 测试 ==========
func testBidirectionalStreaming(ctx context.Context, client pb.TaskServiceClient) {
	fmt.Println("  [4.1] TaskUpdates - 双向流测试")
	stream, err := client.TaskUpdates(ctx)
	if err != nil {
		fmt.Printf("  ❌ 失败: %v\n", err)
		return
	}
	
	var wg sync.WaitGroup
	wg.Add(1)
	
	// 接收协程
	go func() {
		defer wg.Done()
		count := 0
		for {
			resp, err := stream.Recv()
			if err != nil {
				break
			}
			count++
			if resp.ChangeEvent != nil {
				fmt.Printf("  ✅ 收到事件 #%d: %s -> %s\n", 
					count, resp.ChangeEvent.FromStatus, resp.ChangeEvent.ToStatus)
			} else if resp.Task != nil {
				fmt.Printf("  ✅ 收到响应 #%d: ID=%s, Name=%s\n", 
					count, resp.Task.Id, resp.Task.Name)
			}
			if count >= 3 {
				break
			}
		}
	}()
	
	// 发送创建请求
	err = stream.Send(&pb.TaskUpdateRequest{
		RequestId:  "req-1",
		UpdateType: "create",
		Create: &pb.CreateTaskRequest{
			Name: "bidirectional-task-1",
		},
	})
	if err != nil {
		fmt.Printf("  ❌ 发送失败: %v\n", err)
	}
	time.Sleep(100 * time.Millisecond)
	
	// 发送更新请求 (使用已存在的任务)
	if testTaskId != "" {
		err = stream.Send(&pb.TaskUpdateRequest{
			RequestId:  "req-2",
			UpdateType: "update",
			Update: &pb.UpdateTaskRequest{
				Id:     testTaskId,
				Status: pb.TaskStatus(4), // FAILED
			},
		})
		time.Sleep(100 * time.Millisecond)
	}
	
	// 发送未知类型
	err = stream.Send(&pb.TaskUpdateRequest{
		RequestId:  "req-3",
		UpdateType: "unknown",
	})
	time.Sleep(100 * time.Millisecond)
	
	// 关闭发送端
	stream.CloseSend()
	wg.Wait()
	
	fmt.Println("  ✅ Bidirectional Streaming 测试完成")
}
