#!/bin/bash

pushd frontend
npm run serve &
popd

# pushd services/frontend-host
# source ./run.sh &
# popd

pushd services/orders
source ./run.sh &
popd

pushd services/products
source ./run.sh &
popd

pushd services/users
source ./run.sh &
popd

pushd local-gateway
source ./run.sh &
popd