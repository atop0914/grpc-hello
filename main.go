package main

import (
	"os"

	"taskflow/internal/config"
	"taskflow/internal/logger"
	"taskflow/internal/server"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	if err := logger.Init(cfg.Server.EnableDebug); err != nil {
		os.Stderr.WriteString("Failed to initialize logger: " + err.Error())
		os.Exit(1)
	}
	defer logger.Sync()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		logger.Fatalf("Configuration error: %v", err)
	}

	logger.Infof("Starting Task Scheduler Server...")
	logger.Infof("Debug mode: %v", cfg.Server.EnableDebug)
	logger.Infof("HTTP server: %s", cfg.GetHTTPAddr())

	// 创建并启动服务器
	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		logger.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
