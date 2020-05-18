#!/bin/bash

#
# For local debugging / testing only
# ==================================
# Stop everything started by the start-local.sh script
#

# Needed
dapr stop --app-id products
dapr stop --app-id orders
dapr stop --app-id cart
dapr stop --app-id api-gateway
dapr stop --app-id users
docker rm -f api-gateway
pkill -f ".*vue-cli-service.*"

# Probably unnecessary 
pkill dapr
pkill daprd
pkill nginx
pkill frontend-host
pkill orders
pkill cart
pkill products
pkill users
pkill main