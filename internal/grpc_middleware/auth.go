package grpc_middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKeys context keys
type ContextKeys struct {
	UserID   string
	UserName string
	Token    string
}

var contextKeys = &ContextKeys{}

// PublicMethods public methods that don't require authentication
var PublicMethods = map[string]bool{
	"/taskflow.TaskService/HealthCheck": true,
	"/taskflow.TaskService/Login":       true,
}

// AuthConfig authentication config
type AuthConfig struct {
	Secret          string
	TokenExpireHours int
}

// DefaultAuthConfig default auth config
var DefaultAuthConfig = &AuthConfig{
	Secret:          "taskflow-secret-key",
	TokenExpireHours: 24,
}

// UnaryAuthInterceptor creates unary auth interceptor
func UnaryAuthInterceptor(cfg *AuthConfig) grpc.UnaryServerInterceptor {
	if cfg == nil {
		cfg = DefaultAuthConfig
	}
	
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for public methods
		if PublicMethods[info.FullMethod] {
			return handler(ctx, req)
		}
		
		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		
		// Get authorization header
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}
		
		// Parse Bearer token
		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format")
		}
		
		// In production, validate JWT token here
		// For now, extract user info from token (simplified)
		userID, userName, err := validateToken(token, cfg.Secret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		
		// Add user info to context
		ctx = context.WithValue(ctx, contextKeys.UserID, userID)
		ctx = context.WithValue(ctx, contextKeys.UserName, userName)
		ctx = context.WithValue(ctx, contextKeys.Token, token)
		
		return handler(ctx, req)
	}
}

// StreamAuthInterceptor creates stream auth interceptor
func StreamAuthInterceptor(cfg *AuthConfig) grpc.StreamServerInterceptor {
	if cfg == nil {
		cfg = DefaultAuthConfig
	}
	
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip auth for public methods
		if PublicMethods[info.FullMethod] {
			return handler(srv, ss)
		}
		
		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return status.Errorf(codes.Unauthenticated, "missing authorization header")
		}
		
		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if token == authHeader[0] {
			return status.Errorf(codes.Unauthenticated, "invalid authorization format")
		}
		
		userID, userName, err := validateToken(token, cfg.Secret)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		
		// Add user info to context
		ctx := ss.Context()
		ctx = context.WithValue(ctx, contextKeys.UserID, userID)
		ctx = context.WithValue(ctx, contextKeys.UserName, userName)
		ctx = context.WithValue(ctx, contextKeys.Token, token)
		
		wrappedStream := &serverStream{
			ServerStream: ss,
			ctx:         ctx,
		}
		
		return handler(srv, wrappedStream)
	}
}

// validateToken validates token and returns user info
// In production, implement proper JWT validation
func validateToken(token, secret string) (userID, userName string, err error) {
	// Simplified validation - in production use proper JWT library
	// For now, accept any non-empty token and extract user info
	if len(token) == 0 {
		return "", "", status.Errorf(codes.InvalidArgument, "empty token")
	}
	
	// Extract user info from token (simplified)
	// In production: verify JWT signature, check expiration, etc.
	userID = "user-" + token[:min(8, len(token))]
	userName = "User"
	
	return userID, userName, nil
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(contextKeys.UserID).(string); ok {
		return userID
	}
	return ""
}

// GetUserName extracts user name from context
func GetUserName(ctx context.Context) string {
	if userName, ok := ctx.Value(contextKeys.UserName).(string); ok {
		return userName
	}
	return ""
}

// GetToken extracts token from context
func GetToken(ctx context.Context) string {
	if token, ok := ctx.Value(contextKeys.Token).(string); ok {
		return token
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// serverStream wraps grpc.ServerStream to override context
type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}
