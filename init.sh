#!/bin/bash

# Stop and remove containers
docker-compose down -v

# Build Docker images
docker-compose build

# Start all services
docker-compose up -d