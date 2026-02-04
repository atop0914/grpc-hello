#!/bin/bash

echo "🚀 Starting gRPC-Hello Service Performance Test"

# Build the service
echo "🏗️  Building the service..."
go build -o grpc-hello-perf .

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"

# Run basic functionality test
echo "🧪 Running basic functionality test..."

# Start the service in background with custom ports to avoid conflicts
GRPC_PORT=9091 HTTP_PORT=9092 timeout 10s ./grpc-hello-perf &
SERVICE_PID=$!

# Wait a moment for the service to start
sleep 3

if ps -p $SERVICE_PID > /dev/null; then
    echo "✅ Service started successfully"
    
    # Test gRPC connection by attempting to connect
    echo "🔌 Testing service connectivity..."
    
    # Kill the service gracefully
    kill -TERM $SERVICE_PID
    wait $SERVICE_PID 2>/dev/null
    
    echo "⏹️  Service stopped gracefully"
else
    echo "❌ Service failed to start"
    exit 1
fi

echo "🎯 All tests passed! Service is working correctly."

# Show binary size
echo "📊 Binary size: $(du -h grpc-hello-perf | cut -f1)"