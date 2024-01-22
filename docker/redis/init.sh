#!/bin/bash

# Build Docker images
docker-compose build

# Stop and remove containers, networks, images, and volumes
docker-compose down -v

# Start all services
docker-compose up -d --remove-orphans


# Wait for the services to start.
echo "Waiting for services to start..."
sleep 3

docker-compose exec redis-node-1 sh -c "echo 'yes' | redis-cli --cluster create redis-node-1:6379 redis-node-2:6379 redis-node-3:6379 redis-node-4:6379 redis-node-5:6379 redis-node-6:6379 --cluster-replicas 1"