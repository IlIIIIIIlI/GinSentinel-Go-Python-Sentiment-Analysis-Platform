#!/bin/bash

# Create directories if they don't exist
mkdir -p python-service/gen
mkdir -p logs

# Start the services with Docker Compose
echo "Starting sentiment analysis services..."
docker-compose down
docker-compose build --no-cache  # 使用 --no-cache 确保完全重新构建
docker-compose up -d

# Wait for services to start
echo "Waiting for services to start up..."
sleep 10

# Show the running containers
echo "Running containers:"
docker-compose ps

echo "Sentiment service is now available at: http://localhost:8080"
echo "To see the logs, run: docker-compose logs -f"
echo "To stop the services, run: docker-compose down"