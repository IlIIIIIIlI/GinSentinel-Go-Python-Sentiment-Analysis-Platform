FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o bin/sentiment-service cmd/main.go

# Use a smaller image for the runtime
FROM alpine:3.18

WORKDIR /app

# Install necessary dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/bin/sentiment-service .
COPY --from=builder /app/configs /app/configs

# Create directory for logs
RUN mkdir -p /app/logs

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./sentiment-service"]