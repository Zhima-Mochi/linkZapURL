#!/bin/bash

# Build Docker images
docker-compose build

# Stop
docker-compose down

# Start all services
docker-compose up -d