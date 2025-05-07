.PHONY: build run-local stop clean gen-proto test docker-up docker-down docker-logs

# Build the Go application
build:
	go build -o bin/sentiment-service cmd/main.go

# Run locally
run-local: build
	./bin/sentiment-service

# Stop local services
stop:
	pkill -f sentiment-service || true

# Clean build artifacts
clean:
	rm -rf bin/*
	rm -rf internal/gen/*
	rm -rf python-service/gen/*

# Generate Protocol Buffers code
gen-proto:
	# 1. Generate Go and gRPC code with buf
	buf generate
	# 2. Generate Python code manually
	python -m grpc_tools.protoc --python_out=./python-service/gen --grpc_python_out=./python-service/gen -I./proto ./proto/sentiment/v1/sentiment.proto

# Run tests
test:
	go test -v ./...

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Install Buf CLI
install-buf:
	go install github.com/bufbuild/buf/cmd/buf@latest

# Initialize Buf in the project
init-buf:
	buf mod init

# Lint Protocol Buffers
lint-proto:
	buf lint

# Check for breaking changes in Protocol Buffers
breaking-proto:
	buf breaking --against '.git#branch=main'

# Build and push Docker images
docker-build:
	docker-compose build

# Run migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/migration/main.go

# Run locally with Docker dependencies (Postgres and RabbitMQ)
run-deps:
	docker-compose up -d postgres rabbitmq

# Create a development database
create-dev-db:
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to start..."
	@sleep 5
	docker-compose exec postgres psql -U user -d postgres -c "CREATE DATABASE sentiment_dev;"

# Generate example config file
gen-config:
	@echo "Generating example config file..."
	@cp configs/config.yaml configs/config.example.yaml

# Helpful targets for development
help:
	@echo "Available targets:"
	@echo "  build          - Build the Go application"
	@echo "  run-local      - Run the application locally"
	@echo "  stop           - Stop local services"
	@echo "  clean          - Clean build artifacts"
	@echo "  gen-proto      - Generate Protocol Buffers code"
	@echo "  test           - Run tests"
	@echo "  docker-up      - Start all containers"
	@echo "  docker-down    - Stop all containers"
	@echo "  docker-logs    - Show container logs"
	@echo "  install-buf    - Install Buf CLI"
	@echo "  init-buf       - Initialize Buf in the project"
	@echo "  lint-proto     - Lint Protocol Buffers"
	@echo "  breaking-proto - Check for breaking changes in Protocol Buffers"
	@echo "  docker-build   - Build Docker images"
	@echo "  migrate        - Run database migrations"
	@echo "  run-deps       - Run dependencies (Postgres and RabbitMQ)"
	@echo "  create-dev-db  - Create a development database"
	@echo "  gen-config     - Generate example config file"
