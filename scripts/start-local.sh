#!/bin/bash

#
# For local debugging / testing only
# ==================================
# This script starts all services + gateway + frontend
# All will run in hot-reload mode so will restart on code changes
# See run.sh in each directory for details
#

pushd web/frontend
npm run serve &
popd

# pushd cmd/frontend-host
# source ./run.sh &
# popd

pushd cmd/orders
source ./run.sh &
popd

pushd cmd/cart
source ./run.sh &
popd

pushd cmd/products
source ./run.sh &
popd

pushd cmd/users
source ./run.sh &
popd

sleep 5

pushd scripts/local-gateway
source ./run.sh &
popd