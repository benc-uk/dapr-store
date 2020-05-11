#!/bin/bash
#dapr run --app-id cart --app-port 9001 --log-level warn realize start
dapr run --app-id cart --app-port 9001 --log-level warn go run main.go routes.go
