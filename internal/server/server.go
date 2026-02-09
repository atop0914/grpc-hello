package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"grpc-hello/internal/config"
	"grpc-hello/internal/middleware"
)

// Server HTTP服务封装
type Server struct {
	cfg      *config.Config
	httpServer *http.Server
	started  bool
	startMutex sync.Mutex
}

// NewServer 创建服务实例
func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// Start 启动服务
func (s *Server) Start() error {
	s.startMutex.Lock()
	defer s.startMutex.Unlock()

	if s.started {
		return fmt.Errorf("server already started")
	}

	if err := s.startHTTP(); err != nil {
		return fmt.Errorf("failed to start HTTP: %w", err)
	}

	s.started = true
	log.Printf("Server started: HTTP=%s", s.cfg.GetHTTPAddr())

	s.waitForShutdown()

	return nil
}

// startHTTP 启动HTTP服务
func (s *Server) startHTTP() error {
	if s.cfg.Server.EnableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.RequestID(),
		middleware.CORS(),
		middleware.Timeout(s.cfg.GetTimeout()),
	)

	// 简单健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.httpServer = &http.Server{
		Addr:         s.cfg.GetHTTPAddr(),
		Handler:     router,
		ReadTimeout:  s.cfg.GetTimeout(),
		WriteTimeout: s.cfg.GetTimeout(),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("HTTP server listening on %s", s.cfg.GetHTTPAddr())
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}

// waitForShutdown 等待退出信号并优雅关闭
func (s *Server) waitForShutdown() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	log.Println("Shutting down server...")

	gracefulTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTimeout)
	defer cancel()

	s.startMutex.Lock()
	s.started = false
	s.startMutex.Unlock()

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		} else {
			log.Println("Server stopped gracefully")
		}
	}

	log.Println("Server stopped")
}

// GetHTTPAddr 获取HTTP地址
func (s *Server) GetHTTPAddr() string {
	return s.cfg.GetHTTPAddr()
}
