# Application name
APP_NAME = blumfield

# Directories and Files
BUILD_DIR = ./build

# Default target
.DEFAULT_GOAL := build

# Commands
.PHONY: build test db prod run clean sqlc

# Build the Go project
build:
	go build -v -o $(BUILD_DIR)/$(APP_NAME) ./main.go

# Run tests
test:
	go test -v -race -timeout 30s ./...

# Production build for different platforms
prod:
	@if [ "$(filter windows,$(MAKECMDGOALS))" != "" ]; then \
		GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-win-x86.exe -v ./main.go; \
	elif [ "$(filter macos,$(MAKECMDGOALS))" != "" ]; then \
		GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 -v ./main.go; \
	else \
		GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 -v ./main.go; \
	fi

# Run the app
run:
	go run ./main.go

# Clean build artifacts
clean:
	rm -f $(BUILD_DIR)/$(APP_NAME)*
