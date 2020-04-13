#!/bin/bash

dapr stop --app-id products
dapr stop --app-id orders
dapr stop --app-id api-gateway
dapr stop --app-id users
docker rm -f api-gateway
pkill dapr
pkill daprd
pkill nginx
pkill frontend-host
pkill orders
pkill products
pkill users