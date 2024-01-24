#!/bin/bash

# Build Docker images
docker-compose build

# Stop
docker-compose down -v

# Start all services
docker-compose up -d

# Wait for the services to start.
echo "Waiting for services to start..."
sleep 10

# Initialize configuration server
docker-compose exec configsvr01 bash "/scripts/init-configserver.js"

# Initialize each shard
docker-compose exec shard01-a bash "/scripts/init-shard01.js"
docker-compose exec shard02-a bash "/scripts/init-shard02.js"
docker-compose exec shard03-a bash "/scripts/init-shard03.js"

# Wait again to ensure shard initialization is complete
echo "Waiting for shards to initialize..."
sleep 10

# Initialize router
docker-compose exec router01 sh -c "mongosh < /scripts/init-router.js"

# Wait to ensure the router is initialized
echo "Waiting for router to initialize..."
sleep 10

# Apply authentication settings to the config server, each shard, and the router
docker-compose exec configsvr01 bash "/scripts/auth.js"
docker-compose exec shard01-a bash "/scripts/auth.js"
docker-compose exec shard02-a bash "/scripts/auth.js"
docker-compose exec shard03-a bash "/scripts/auth.js"

# Wait to ensure authentication setup is complete
echo "Waiting for authentication setup..."
sleep 10

# Execute the final initialization script
docker-compose exec router01 sh -c "mongosh -u root --authenticationDatabase admin -p password < /scripts/mongo-init.js"

echo "MongoDB Sharded Cluster is ready!"