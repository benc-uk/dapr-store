#!/bin/bash
helm install api-gateway stable/nginx-ingress -f chart-values.yaml
