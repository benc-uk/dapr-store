#!/bin/bash

#
# This script starts the local version of the API gateway
# Wraps nginx-proxy.sh in `dapr run` so that nginx runs Dapr-ized
#

scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
dapr run --app-id api-gateway --app-port 9000 --log-level warn "$scriptDir"/nginx-proxy.sh --components-path "$scriptDir"/../../components