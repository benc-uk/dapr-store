#!/bin/bash

# pushd cmd/frontend-host
# source ./run.sh &
# popd
pushd web/frontend
npm run serve &
popd

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