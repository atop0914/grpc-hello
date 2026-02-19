package grpc_middleware

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// LoggerConfig logger config
type LoggerConfig struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

// defaultLoggerConfig default config
var defaultLoggerConfig = &LoggerConfig{
	InfoLogger:  log.Default(),
	ErrorLogger: log.Default(),
}

// RequestIDHeader request ID header name
const RequestIDHeader = "x-request-id"

// UnaryLoggerInterceptor creates unary logger interceptor
func UnaryLoggerInterceptor(cfg *LoggerConfig) grpc.UnaryServerInterceptor {
	if cfg == nil {
		cfg = defaultLoggerConfig
	}
	
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		
		// Generate or extract request ID
		requestID := generateRequestID(ctx)
		ctx = context.WithValue(ctx, "request_id", requestID)
		
		// Log request
		cfg.InfoLogger.Printf("[%s] RPC started: %s", requestID, info.FullMethod)
		
		// Get metadata
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			cfg.InfoLogger.Printf("[%s] Metadata: %v", requestID, md)
		}
		
		// Call handler
		resp, err := handler(ctx, req)
		
		// Log response
		duration := time.Since(startTime)
		if err != nil {
			cfg.InfoLogger.Printf("[%s] RPC failed: %s, duration: %v, error: %v", 
				requestID, info.FullMethod, duration, err)
		} else {
			cfg.InfoLogger.Printf("[%s] RPC completed: %s, duration: %v", 
				requestID, info.FullMethod, duration)
		}
		
		return resp, err
	}
}

// StreamLoggerInterceptor creates stream logger interceptor
func StreamLoggerInterceptor(cfg *LoggerConfig) grpc.StreamServerInterceptor {
	if cfg == nil {
		cfg = defaultLoggerConfig
	}
	
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()
		
		// Generate or extract request ID
		requestID := generateRequestID(ss.Context())
		ctx := context.WithValue(ss.Context(), "request_id", requestID)
		
		// Log stream start
		cfg.InfoLogger.Printf("[%s] Stream started: %s", requestID, info.FullMethod)
		
		// Create wrapped stream
		wrappedStream := &loggingStream{
			ServerStream: ss,
			ctx:          ctx,
			requestID:    requestID,
			config:       cfg,
			method:      info.FullMethod,
			startTime:   startTime,
		}
		
		// Call handler
		err := handler(srv, wrappedStream)
		
		// Log stream end
		duration := time.Since(startTime)
		if err != nil {
			cfg.InfoLogger.Printf("[%s] Stream failed: %s, duration: %v, error: %v", 
				requestID, info.FullMethod, duration, err)
		} else {
			cfg.InfoLogger.Printf("[%s] Stream completed: %s, duration: %v", 
				requestID, info.FullMethod, duration)
		}
		
		return err
	}
}

// loggingStream wraps grpc.ServerStream for logging
type loggingStream struct {
	grpc.ServerStream
	ctx        context.Context
	requestID  string
	config     *LoggerConfig
	method     string
	startTime  time.Time
}

func (s *loggingStream) Context() context.Context {
	return s.ctx
}

func (s *loggingStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err != nil {
		s.config.InfoLogger.Printf("[%s] Stream send error: %v", s.requestID, err)
	}
	return err
}

func (s *loggingStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err != nil && err != io.EOF {
		s.config.InfoLogger.Printf("[%s] Stream receive error: %v", s.requestID, err)
	}
	return err
}

// ========== Recovery Interceptor ==========

// UnaryRecoveryInterceptor creates unary recovery interceptor
func UnaryRecoveryInterceptor(cfg *LoggerConfig) grpc.UnaryServerInterceptor {
	if cfg == nil {
		cfg = defaultLoggerConfig
	}
	
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				cfg.ErrorLogger.Printf("[PANIC] Recovered in unary RPC: %s, error: %v", 
					info.FullMethod, r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		
		return handler(ctx, req)
	}
}

// StreamRecoveryInterceptor creates stream recovery interceptor
func StreamRecoveryInterceptor(cfg *LoggerConfig) grpc.StreamServerInterceptor {
	if cfg == nil {
		cfg = defaultLoggerConfig
	}
	
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				cfg.ErrorLogger.Printf("[PANIC] Recovered in stream RPC: %s, error: %v", 
					info.FullMethod, r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		
		return handler(srv, ss)
	}
}

// ========== Utility Functions ==========

// generateRequestID generates or extracts request ID
func generateRequestID(ctx context.Context) string {
	// Try to get from context
	if rid, ok := ctx.Value("request_id").(string); ok && rid != "" {
		return rid
	}
	
	// Try to get from metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get(RequestIDHeader); len(ids) > 0 && ids[0] != "" {
			return ids[0]
		}
	}
	
	// Generate new ID
	return generateID()
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if rid, ok := ctx.Value("request_id").(string); ok {
		return rid
	}
	return ""
}
