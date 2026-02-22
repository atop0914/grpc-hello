package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TaskCount - task count gauge
	TaskCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "taskflow_tasks_total",
		Help: "Total number of tasks by status",
	}, []string{"status"})

	// TaskDuration - task execution duration
	TaskDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "taskflow_task_duration_seconds",
		Help:    "Task execution duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"task_type", "status"})

	// TaskErrors - task error counter
	TaskErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "taskflow_task_errors_total",
		Help: "Total number of task errors",
	}, []string{"task_type", "error_type"})

	// SchedulerDelay - scheduler delay histogram
	SchedulerDelay = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "taskflow_scheduler_delay_seconds",
		Help:    "Scheduler delay from pending to running",
		Buckets: prometheus.DefBuckets,
	})

	// GRPCRequests - gRPC request counter
	GRPCRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "taskflow_grpc_requests_total",
		Help: "Total number of gRPC requests",
	}, []string{"method", "status"})

	// GRPCLatency - gRPC latency histogram
	GRPCLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "taskflow_grpc_latency_seconds",
		Help:    "gRPC request latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})
)

// RecordTaskStatus records task status count
func RecordTaskStatus(status string, count int) {
	TaskCount.WithLabelValues(status).Set(float64(count))
}

// RecordTaskDuration records task execution duration
func RecordTaskDuration(taskType, status string, duration float64) {
	TaskDuration.WithLabelValues(taskType, status).Observe(duration)
}

// RecordTaskError records task error
func RecordTaskError(taskType, errorType string) {
	TaskErrors.WithLabelValues(taskType, errorType).Inc()
}

// RecordSchedulerDelay records scheduler delay
func RecordSchedulerDelay(delay float64) {
	SchedulerDelay.Observe(delay)
}

// RecordGRPCRequest records gRPC request
func RecordGRPCRequest(method, status string) {
	GRPCRequests.WithLabelValues(method, status).Inc()
}

// RecordGRPCLatency records gRPC latency
func RecordGRPCLatency(method string, duration float64) {
	GRPCLatency.WithLabelValues(method).Observe(duration)
}
