#!/bin/bash

pushd services/frontend-host
source ./run.sh &
popd

pushd services/orders
source ./run.sh &
popd

pushd services/cart
source ./run.sh &
popd

pushd services/products
source ./run.sh &
popd

pushd services/users
source ./run.sh &
popd

sleep 5

pushd local-gateway
source ./run.sh &
popd