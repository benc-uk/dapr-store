#!/bin/bash

rm ./frontend-host
export STATIC_DIR=../../frontend/dist 
go build -o frontend-host
./frontend-host
