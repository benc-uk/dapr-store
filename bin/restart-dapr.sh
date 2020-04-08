#!/bin/bash

kubectl rollout restart -n dapr-system deploy/dapr-operator
kubectl rollout restart -n dapr-system deploy/dapr-placement
kubectl rollout restart -n dapr-system deploy/dapr-sentry
kubectl rollout restart -n dapr-system deploy/dapr-sidecar-injector
