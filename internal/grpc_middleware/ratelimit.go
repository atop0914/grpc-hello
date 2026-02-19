package grpc_middleware

import (
	"context"
	"math"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimiterConfig rate limiter config
type RateLimiterConfig struct {
	RequestsPerSecond float64       // requests per second
	BurstSize        int           // max burst size
	ClientKeyFunc    func(ctx context.Context) string // function to get client key
}

// defaultRateLimiterConfig default config
var defaultRateLimiterConfig = &RateLimiterConfig{
	RequestsPerSecond: 100,
	BurstSize:        200,
	ClientKeyFunc:    defaultClientKeyFunc,
}

// TokenBucketLimiter token bucket rate limiter
type TokenBucketLimiter struct {
	mu         sync.RWMutex
	tokens     map[string]*bucket
	config     *RateLimiterConfig
}

type bucket struct {
	tokens     float64
	maxTokens  float64
	lastUpdate time.Time
	refillRate float64
}

// NewTokenBucketLimiter creates new token bucket limiter
func NewTokenBucketLimiter(cfg *RateLimiterConfig) *TokenBucketLimiter {
	if cfg == nil {
		cfg = defaultRateLimiterConfig
	}
	if cfg.ClientKeyFunc == nil {
		cfg.ClientKeyFunc = defaultClientKeyFunc
	}
	
	return &TokenBucketLimiter{
		tokens: make(map[string]*bucket),
		config: cfg,
	}
}

// defaultClientKeyFunc default client key function
func defaultClientKeyFunc(ctx context.Context) string {
	// Use user ID if available, otherwise use "anonymous"
	if userID := GetUserID(ctx); userID != "" {
		return userID
	}
	return "anonymous"
}

// allow checks if request is allowed
func (r *TokenBucketLimiter) allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	b, exists := r.tokens[key]
	
	if !exists {
		// Create new bucket
		r.tokens[key] = &bucket{
			tokens:     float64(r.config.BurstSize),
			maxTokens:  float64(r.config.BurstSize),
			lastUpdate: now,
			refillRate: r.config.RequestsPerSecond,
		}
		return true
	}
	
	// Calculate elapsed time and refill tokens
	elapsed := now.Sub(b.lastUpdate).Seconds()
	b.tokens = math.Min(b.maxTokens, b.tokens+elapsed*b.refillRate)
	b.lastUpdate = now
	
	// Check if enough tokens
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	
	return false
}

// UnaryRateLimiter creates unary rate limiter interceptor
func UnaryRateLimiter(limiter *TokenBucketLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		clientKey := limiter.config.ClientKeyFunc(ctx)
		
		if !limiter.allow(clientKey) {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}
		
		return handler(ctx, req)
	}
}

// StreamRateLimiter creates stream rate limiter interceptor
func StreamRateLimiter(limiter *TokenBucketLimiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientKey := limiter.config.ClientKeyFunc(ss.Context())
		
		if !limiter.allow(clientKey) {
			return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}
		
		return handler(srv, ss)
	}
}

// ========== Sliding Window Rate Limiter ==========

// SlidingWindowLimiter sliding window rate limiter
type SlidingWindowLimiter struct {
	mu           sync.RWMutex
	requests     map[string][]time.Time
	maxRequests  int
	windowSize   time.Duration
	config       *RateLimiterConfig
}

// NewSlidingWindowLimiter creates new sliding window limiter
func NewSlidingWindowLimiter(maxRequests int, windowSize time.Duration, cfg *RateLimiterConfig) *SlidingWindowLimiter {
	if cfg == nil {
		cfg = defaultRateLimiterConfig
	}
	if cfg.ClientKeyFunc == nil {
		cfg.ClientKeyFunc = defaultClientKeyFunc
	}
	
	return &SlidingWindowLimiter{
		requests:   make(map[string][]time.Time),
		maxRequests: maxRequests,
		windowSize: windowSize,
		config:     cfg,
	}
}

// allowSliding checks if request is allowed using sliding window
func (r *SlidingWindowLimiter) allowSliding(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	windowStart := now.Add(-r.windowSize)
	
	// Get existing requests
	times, exists := r.requests[key]
	if !exists {
		r.requests[key] = []time.Time{now}
		return true
	}
	
	// Filter out old requests
	validTimes := make([]time.Time, 0)
	for _, t := range times {
		if t.After(windowStart) {
			validTimes = append(validTimes, t)
		}
	}
	
	// Check limit
	if len(validTimes) >= r.maxRequests {
		r.requests[key] = validTimes
		return false
	}
	
	// Add new request
	r.requests[key] = append(validTimes, now)
	return true
}

// UnarySlidingRateLimiter creates unary sliding window rate limiter
func UnarySlidingRateLimiter(limiter *SlidingWindowLimiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		clientKey := limiter.config.ClientKeyFunc(ctx)
		
		if !limiter.allowSliding(clientKey) {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}
		
		return handler(ctx, req)
	}
}

// StreamSlidingRateLimiter creates stream sliding window rate limiter
func StreamSlidingRateLimiter(limiter *SlidingWindowLimiter) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		clientKey := limiter.config.ClientKeyFunc(ss.Context())
		
		if !limiter.allowSliding(clientKey) {
			return status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}
		
		return handler(srv, ss)
	}
}
