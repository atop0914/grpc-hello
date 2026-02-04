package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	helloworldpb "grpc-hello/proto/helloworld"
	"grpc-hello/route"
	"grpc-hello/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworldpb.UnimplementedGreeterServer
	// Statistics fields
	mu            sync.RWMutex
	totalRequests int64
	uniqueNames   map[string]int64
	nameFrequency map[string]int64  // Track frequency of each name
	lastRequest   time.Time
	config        *config.Config
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) *server {
	return &server{
		uniqueNames:   make(map[string]int64),
		nameFrequency: make(map[string]int64),
		config:        cfg,
	}
}

// updateStats updates the statistics for a given name
func (s *server) updateStats(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.totalRequests++
	if name != "" {
		lowerName := strings.ToLower(name)
		s.uniqueNames[lowerName]++
		s.nameFrequency[lowerName]++
	}
	s.lastRequest = time.Now()
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworldpb.HelloRequest) (*helloworldpb.HelloReply, error) {
	// Validate input
	if in.GetNameTest() == "" {
		in.NameTest = "World" // Default value if empty
	}
	
	// Update statistics
	s.updateStats(in.GetNameTest())
	
	// Determine greeting based on language
	greeting := getGreetingByLanguage(in.Language)
	
	reply := &helloworldpb.HelloReply{
		TestMessage: fmt.Sprintf("%s %s!", greeting, in.GetNameTest()),
		Timestamp:   time.Now().Unix(),
		Language:    in.Language,
		Tags:        in.Tags,
	}
	
	log.Printf("SayHello called: %s in %s -> %s", in.GetNameTest(), in.Language, reply.TestMessage)
	
	return reply, nil
}

// SayHelloMultiple implements multiple greetings
func (s *server) SayHelloMultiple(ctx context.Context, in *helloworldpb.HelloMultipleRequest) (*helloworldpb.HelloMultipleReply, error) {
	var greetings []*helloworldpb.HelloReply
	
	maxGreetings := s.config.Features.MaxGreetings
	if len(in.Names) > maxGreetings {
		return nil, fmt.Errorf("too many names, maximum allowed: %d", maxGreetings)
	}
	
	for _, name := range in.Names {
		// Update statistics for each name
		s.updateStats(name)
		
		greeting := getGreetingByLanguage("") // Default to English for multiple greetings
		
		reply := &helloworldpb.HelloReply{
			TestMessage: fmt.Sprintf("%s %s! %s", greeting, name, in.CommonMessage),
			Timestamp:   time.Now().Unix(),
		}
		greetings = append(greetings, reply)
	}
	
	reply := &helloworldpb.HelloMultipleReply{
		Greetings:   greetings,
		TotalCount:  int32(len(greetings)),
	}
	
	log.Printf("SayHelloMultiple called with %d names", len(in.Names))
	
	return reply, nil
}

// GetGreetingStats implements statistics retrieval
func (s *server) GetGreetingStats(ctx context.Context, in *helloworldpb.GreetingStatsRequest) (*helloworldpb.GreetingStatsReply, error) {
	if !s.config.Features.EnableStats {
		return nil, fmt.Errorf("statistics feature is disabled")
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	nameFrequency := make(map[string]int32)
	filter := strings.ToLower(in.NameFilter)
	
	for name, count := range s.nameFrequency {
		if filter == "" || strings.Contains(name, filter) {
			nameFrequency[name] = int32(count)
		}
	}
	
	// Sort names by frequency (top 10)
	type kv struct {
		Key   string
		Value int32
	}
	var ss []kv
	for k, v := range nameFrequency {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	
	// Limit to top 10
	if len(ss) > 10 {
		ss = ss[:10]
		// Create a new map with only top 10 items
		newMap := make(map[string]int32)
		for _, kv := range ss {
			newMap[kv.Key] = kv.Value
		}
		nameFrequency = newMap
	}
	
	reply := &helloworldpb.GreetingStatsReply{
		TotalRequests:   int32(s.totalRequests),
		UniqueNames:     int32(len(s.uniqueNames)),
		NameFrequency:   nameFrequency,
		LastRequestTime: s.lastRequest.Unix(),
	}
	
	log.Printf("GetGreetingStats called, returning stats for %d unique names", len(s.uniqueNames))
	
	return reply, nil
}

// getGreetingByLanguage returns the appropriate greeting based on the language code
func getGreetingByLanguage(language string) string {
	if language == "" {
		return "Hello"
	}
	
	switch strings.ToLower(language) {
	case "zh", "chinese":
		return "你好"
	case "es", "spanish":
		return "Hola"
	case "fr", "french":
		return "Bonjour"
	case "ja", "japanese":
		return "こんにちは"
	case "ko", "korean":
		return "안녕하세요"
	case "ru", "russian":
		return "Привет"
	case "de", "german":
		return "Hallo"
	case "it", "italian":
		return "Ciao"
	default:
		return "Hello"
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Set Gin mode based on debug flag
	if cfg.Server.EnableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer()
	helloworldpb.RegisterGreeterServer(grpcServer, NewServer(cfg))
	
	// Enable reflection for debugging tools (only if enabled in config)
	if cfg.Features.EnableReflection {
		reflection.Register(grpcServer)
	}

	// Create gRPC listener
	grpcAddr := fmt.Sprintf(":%s", cfg.Server.GRPCPort)
	grpcLis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on %s", grpcAddr)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Wait a bit for gRPC server to start
	time.Sleep(100 * time.Millisecond)

	// Create a client connection to the gRPC server
	grpcEndpoint := fmt.Sprintf("127.0.0.1:%s", cfg.Server.GRPCPort)
	conn, err := grpc.DialContext(
		context.Background(),
		grpcEndpoint,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(cfg.Server.Timeout),
	)
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	// Create gRPC-Gateway mux
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// Register gRPC-Gateway handlers
	err = helloworldpb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	// Create Gin router
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Mount gRPC-Gateway under /rpc/v1
	router.Any("/rpc/v1/*any", gin.WrapF(gwmux.ServeHTTP))
	
	// Initialize custom routes
	route.InitRoute(router)

	// Graceful shutdown handling
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server using endless
	httpAddr := fmt.Sprintf(":%s", cfg.Server.HTTPPort)
	log.Printf("Starting HTTP server on %s", httpAddr)
	err = endless.ListenAndServe(httpAddr, router)
	if err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}

	// Attempt graceful shutdown
	<-stopCh
	log.Println("Shutting down servers...")
	grpcServer.GracefulStop()
	log.Println("Servers stopped")
}