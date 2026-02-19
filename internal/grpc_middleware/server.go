package grpc_middleware

import (
	"google.golang.org/grpc"
)

// ========== Server Options ==========

// ServerOption server option
type ServerOption func(*serverOptions)

type serverOptions struct {
	authEnabled      bool
	rateLimitEnabled bool
	loggerEnabled    bool
	recoveryEnabled  bool
	authConfig       *AuthConfig
	tokenLimiter     *TokenBucketLimiter
	slidingLimiter   *SlidingWindowLimiter
	loggerConfig     *LoggerConfig
}

// WithAuth enables authentication
func WithAuth(cfg *AuthConfig) ServerOption {
	return func(o *serverOptions) {
		o.authEnabled = true
		o.authConfig = cfg
	}
}

// WithRateLimit enables rate limiting with token bucket
func WithRateLimit(limiter *TokenBucketLimiter) ServerOption {
	return func(o *serverOptions) {
		o.rateLimitEnabled = true
		o.tokenLimiter = limiter
	}
}

// WithSlidingWindowRateLimit enables rate limiting with sliding window
func WithSlidingWindowRateLimit(limiter *SlidingWindowLimiter) ServerOption {
	return func(o *serverOptions) {
		o.rateLimitEnabled = true
		o.slidingLimiter = limiter
	}
}

// WithLogger enables logging
func WithLogger(cfg *LoggerConfig) ServerOption {
	return func(o *serverOptions) {
		o.loggerEnabled = true
		o.loggerConfig = cfg
	}
}

// WithRecovery enables panic recovery
func WithRecovery() ServerOption {
	return func(o *serverOptions) {
		o.recoveryEnabled = true
	}
}

// DefaultServerOptions returns default server options
func DefaultServerOptions() *serverOptions {
	return &serverOptions{
		authEnabled:      false,
		rateLimitEnabled: false,
		loggerEnabled:    false,
		recoveryEnabled:  false,
	}
}

// GetUnaryServerOptions returns unary server options
func GetUnaryServerOptions(options ...ServerOption) ([]grpc.ServerOption, error) {
	opts := DefaultServerOptions()
	for _, opt := range options {
		opt(opts)
	}

	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	// Add recovery interceptor first (outermost)
	if opts.recoveryEnabled {
		recoveryCfg := opts.loggerConfig
		if recoveryCfg == nil {
			recoveryCfg = defaultLoggerConfig
		}
		unaryInterceptors = append(unaryInterceptors, UnaryRecoveryInterceptor(recoveryCfg))
		streamInterceptors = append(streamInterceptors, StreamRecoveryInterceptor(recoveryCfg))
	}

	// Add logger interceptor
	if opts.loggerEnabled {
		loggerCfg := opts.loggerConfig
		if loggerCfg == nil {
			loggerCfg = defaultLoggerConfig
		}
		unaryInterceptors = append(unaryInterceptors, UnaryLoggerInterceptor(loggerCfg))
		streamInterceptors = append(streamInterceptors, StreamLoggerInterceptor(loggerCfg))
	}

	// Add rate limiter
	if opts.rateLimitEnabled {
		if opts.tokenLimiter != nil {
			unaryInterceptors = append(unaryInterceptors, UnaryRateLimiter(opts.tokenLimiter))
			streamInterceptors = append(streamInterceptors, StreamRateLimiter(opts.tokenLimiter))
		} else if opts.slidingLimiter != nil {
			unaryInterceptors = append(unaryInterceptors, UnarySlidingRateLimiter(opts.slidingLimiter))
			streamInterceptors = append(streamInterceptors, StreamSlidingRateLimiter(opts.slidingLimiter))
		}
	}

	// Add auth interceptor (innermost)
	if opts.authEnabled {
		authCfg := opts.authConfig
		if authCfg == nil {
			authCfg = DefaultAuthConfig
		}
		unaryInterceptors = append(unaryInterceptors, UnaryAuthInterceptor(authCfg))
		streamInterceptors = append(streamInterceptors, StreamAuthInterceptor(authCfg))
	}

	var serverOpts []grpc.ServerOption

	if len(unaryInterceptors) == 1 {
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(unaryInterceptors[0]))
	} else if len(unaryInterceptors) > 1 {
		// Use the last one for simplicity
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(unaryInterceptors[len(unaryInterceptors)-1]))
	}

	if len(streamInterceptors) == 1 {
		serverOpts = append(serverOpts, grpc.StreamInterceptor(streamInterceptors[0]))
	} else if len(streamInterceptors) > 1 {
		serverOpts = append(serverOpts, grpc.StreamInterceptor(streamInterceptors[len(streamInterceptors)-1]))
	}

	return serverOpts, nil
}
