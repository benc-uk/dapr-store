#!/bin/bash

rm ./frontend-host
export STATIC_DIR=../../web/frontend/dist 
go build -o frontend-host
./frontend-host
