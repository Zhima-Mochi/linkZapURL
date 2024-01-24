#!/bin/bash

# Build Docker images
docker-compose build

# Stop and remove containers, networks, images, and volumes
docker-compose down -v

# Start all services
docker-compose up -d

# Wait for the services to start.
echo "Waiting for services to start..."
sleep 1

# Initialize Redis cluster
docker-compose exec redis-node-1 bash -c "$(cat ./scripts/init-redis-cluster.sh)"

echo "Redis Cluster is ready!"
