#!/bin/bash

pushd services/api-gateway
./run.sh &
popd

pushd services/orders
./run.sh &
popd