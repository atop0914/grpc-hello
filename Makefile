# Makefile for grpc-hello project

# Build the project
.PHONY: build
build:
	CGO_ENABLED=0 go build -o grpc-hello .

# Run the project
.PHONY: run
run:
	go run main.go

# Install dependencies
.PHONY: deps
deps:
	go mod tidy

# Generate protobuf files (if needed)
.PHONY: proto-gen
proto-gen:
	protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative \
		--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=proto --grpc-gateway_opt=paths=source_relative \
		proto/helloworld/hello_world.proto

# Clean build artifacts
.PHONY: clean
clean:
	rm -f grpc-hello

# Test the project
.PHONY: test
test:
	go test ./...

# Build for different platforms
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o grpc-hello-linux .

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o grpc-hello-mac .

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o grpc-hello-windows.exe .

# All builds
.PHONY: build-all
build-all: build-linux build-mac build-windows