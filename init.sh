#!/bin/bash
# Stop and remove containers
docker-compose -f docker-compose.windows.yml down -v

# Build Docker images
docker-compose -f docker-compose.windows.yml build

# Start all services
docker-compose -f docker-compose.windows.yml up