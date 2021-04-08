#!/bin/bash
dapr run --app-id users --app-port 9003 --log-level warn go run main.go routes.go
#realize start
