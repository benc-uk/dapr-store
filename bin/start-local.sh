#!/bin/bash

pushd local-gateway
source ./run.sh &
popd

pushd services/frontend-host
source ./run.sh &
popd

pushd services/orders
source ./run.sh &
popd