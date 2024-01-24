#!/bin/bash

# Check if Redis cluster is already initialized
if redis-cli -p 7000 -a password cluster info &>/dev/null; then
    echo "Redis cluster is already initialized."
    exit 1
fi

HOST_IP=host.docker.internal

# Start Redis cluster nodes
echo "yes" | redis-cli -h ${HOST_IP} -p 7000 --cluster create \
    ${HOST_IP}:7000 ${HOST_IP}:7001 ${HOST_IP}:7002 \
    ${HOST_IP}:7003 ${HOST_IP}:7004 ${HOST_IP}:7005 \
    --cluster-replicas 1 -a password

# Verify cluster status
redis-cli -p 7000 -a password cluster info
