package main

import (
	"log"
	"os"

	"taskflow/internal/config"
	"taskflow/internal/server"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	log.Printf("Starting Task Scheduler Server...")
	log.Printf("Debug mode: %v", cfg.Server.EnableDebug)
	log.Printf("HTTP server: %s", cfg.GetHTTPAddr())

	// 创建并启动服务器
	srv := server.NewServer(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}
