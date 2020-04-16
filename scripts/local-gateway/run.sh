#!/bin/bash
dapr run --app-id api-gateway --app-port 9000 --log-level error ./nginx-proxy.sh