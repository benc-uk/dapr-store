#!/bin/bash
set -e

ACR_NAME="bcdemo"

docker build . -f build/service.Dockerfile \
--build-arg serviceName=api-gateway \
--build-arg servicePort=9000 \
-t $ACR_NAME.azurecr.io/dapr-store/api-gateway

docker build . -f build/service.Dockerfile \
--build-arg serviceName=orders \
--build-arg servicePort=9001 \
-t $ACR_NAME.azurecr.io/dapr-store/orders

docker push $ACR_NAME.azurecr.io/dapr-store/api-gateway
docker push $ACR_NAME.azurecr.io/dapr-store/orders