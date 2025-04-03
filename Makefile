.PHONY: build run clean install test

# Build variables
BINARY_NAME=geq
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Default target
all: build

# Build the application
build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o ${BINARY_NAME}

# Install the application
install:
	@echo "Installing ${BINARY_NAME}..."
	@go install ${LDFLAGS}

# Run the application
run:
	@go run ${LDFLAGS} main.go $(ARGS)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Get dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download