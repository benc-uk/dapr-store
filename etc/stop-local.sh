#!/bin/bash

dapr stop --app-id orders
dapr stop --app-id api-gateway
pkill orders
pkill api-gateway
