#!/bin/bash

# Build Docker images
docker-compose build

# Stop and remove containers, networks, images, and volumes
docker-compose down -v

# Start all services
docker-compose up -d

echo "Redis Cluster is ready!"
