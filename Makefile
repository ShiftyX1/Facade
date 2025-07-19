
.PHONY: build run test clean help

build:
	mkdir -p build
	go build -o build/facade cmd/facade/main.go

run: build
	./build/facade -c configs/example.yaml -v

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
	@echo "  build     - Build the application"
	@echo "  run       - Build and run with example config"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Install dependencies"
	@echo "  build-all - Cross-compile for multiple platforms"
	@echo "  help      - Show this help"

all: build
