package logger

import (
	"go.uber.org/zap"
)

var (
	// Logger 全局日志实例
	Logger *zap.SugaredLogger
)

// Init 初始化日志
func Init(debug bool) error {
	var cfg zap.Config
	if debug {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	Logger = l.Sugar()
	return nil
}

// InitWithConfig 使用配置初始化
func InitWithConfig(cfg zap.Config) error {
	l, err := cfg.Build()
	if err != nil {
		return err
	}

	Logger = l.Sugar()
	return nil
}

// Sync 同步日志缓冲
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Debug 调试日志
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}

// With 添加上下文
func With(args ...interface{}) *zap.SugaredLogger {
	return Logger.With(args...)
}
