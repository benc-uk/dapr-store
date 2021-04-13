#!/bin/bash

echo "🠶🠶🠶  WARNING! 💥 This will clear all Redis state"
read -r -p "🠶🠶🠶  Press enter to continue, or ctrl-c to exit"

docker info > /dev/null 2>&1 || { echo "Docker is not running!"; exit 1; }

docker run --rm --network host redis redis-cli --scan --pattern "orders||*" | xargs redis-cli del
docker run --rm --network host redis redis-cli --scan --pattern "users||*" | xargs redis-cli del
docker run --rm --network host redis redis-cli --scan --pattern "carts||*" | xargs redis-cli del
