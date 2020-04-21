#!/bin/bash

#
# This script starts the local version of the API gateway
#

dapr run --app-id api-gateway --app-port 9000 --log-level error ./nginx-proxy.sh