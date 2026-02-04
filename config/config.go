package config

import (
	"os"
	"strconv"
	"time"
)

// ServerConfig holds server-related configurations
type ServerConfig struct {
	GRPCPort     string
	HTTPPort     string
	EnableDebug  bool
	Timeout      time.Duration
	MaxConns     int
	LogLevel     string
}

// FeatureFlags holds feature toggle configurations
type FeatureFlags struct {
	EnableReflection bool
	EnableStats      bool
	EnableMetrics    bool
	MaxGreetings     int
}

// Config represents the complete configuration
type Config struct {
	Server   ServerConfig
	Features FeatureFlags
}

// LoadConfig loads configuration from environment variables and defaults
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			GRPCPort:    getEnvOrDefault("GRPC_PORT", "8080"),
			HTTPPort:    getEnvOrDefault("HTTP_PORT", "8090"),
			EnableDebug: getEnvBoolOrDefault("ENABLE_DEBUG", false),
			Timeout:     time.Duration(getEnvIntOrDefault("SERVER_TIMEOUT", 30)) * time.Second,
			MaxConns:    getEnvIntOrDefault("MAX_CONNECTIONS", 1000),
			LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
		},
		Features: FeatureFlags{
			EnableReflection: getEnvBoolOrDefault("ENABLE_REFLECTION", false), // Match debug mode by default
			EnableStats:      getEnvBoolOrDefault("ENABLE_STATS", true),
			EnableMetrics:    getEnvBoolOrDefault("METRICS_ENABLED", true),
			MaxGreetings:     getEnvIntOrDefault("MAX_GREETINGS", 100),
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// Helper functions for environment variable parsing
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		// Parse boolean: "true", "1", "yes", "on" => true; others => false
		switch value {
		case "true", "1", "yes", "on":
			return true
		default:
			return false
		}
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}