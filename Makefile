BINARY_NAME=zerotier-webhook
DOCKER_IMAGE_NAME=$(BINARY_NAME):latest

.PHONY: all build run test clean docker-build docker-run help

all: build test

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) events.db

build:
	@echo "Building $(BINARY_NAME)..."
	@CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o $(BINARY_NAME) ./cmd/server

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME) .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE_NAME)

help:
	@echo "Usage: make [TARGET]"
	@echo ""
	@echo "Targets:"
	@echo "  all           Build and run tests"
	@echo "  build         Build the application"
	@echo "  run           Build and run the application"
	@echo "  test          Run tests"
	@echo "  clean         Remove binary and database"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Build and run Docker container"
	@echo "  help          Show this help message"
