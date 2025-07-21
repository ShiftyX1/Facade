
.PHONY: build build-with-version run test clean help deps build-all

# Detect OS for cross-platform compatibility
ifeq ($(OS),Windows_NT)
    BINARY_NAME = facade.exe
    BUILD_PATH = build/facade.exe
else
    BINARY_NAME = facade
    BUILD_PATH = build/facade
endif

build:
	mkdir -p build
	go build \
		-ldflags="-s -w -X 'main.Version=dev-$(shell git rev-parse --short HEAD)' -X 'main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" \
		-o $(BUILD_PATH) \
		cmd/facade/main.go

build-with-version:
	mkdir -p build
	go build \
		-ldflags="-s -w -X 'main.Version=$(shell git describe --tags --always --dirty)' -X 'main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" \
		-o $(BUILD_PATH) \
		cmd/facade/main.go

run: build
	$(BUILD_PATH) -c configs/example.yaml -v

test: build
	go test ./...

clean:
	rm -rf build/

deps:
	go mod tidy
	go mod download

build-all: clean
	GOOS=linux GOARCH=amd64 go build -o facade-linux-amd64 cmd/facade/main.go
	GOOS=windows GOARCH=amd64 go build -o facade-windows-amd64.exe cmd/facade/main.go
	GOOS=darwin GOARCH=amd64 go build -o facade-darwin-amd64 cmd/facade/main.go
	GOOS=darwin GOARCH=arm64 go build -o facade-darwin-arm64 cmd/facade/main.go

help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  build-with-version - Build with git version and build time"
	@echo "  run                - Build and run with example config"
	@echo "  test               - Run tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  deps               - Install dependencies"
	@echo "  build-all          - Cross-compile for multiple platforms"
	@echo "  help               - Show this help"

all: build
