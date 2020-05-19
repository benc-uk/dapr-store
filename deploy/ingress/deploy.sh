#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

helm upgrade --install api-gateway stable/nginx-ingress -f $DIR/chart-values.yaml
kubectl apply -f $DIR/rules.yaml
