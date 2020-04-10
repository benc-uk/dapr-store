#!/bin/bash
set -e

ACR_NAME="bcdemo"

docker build . -f build/service.Dockerfile \
--build-arg serviceName=products \
--build-arg servicePort=9002 \
-t $ACR_NAME.azurecr.io/dapr-store/products

docker build . -f build/service.Dockerfile \
--build-arg serviceName=orders \
--build-arg servicePort=9001 \
-t $ACR_NAME.azurecr.io/dapr-store/orders

docker build . -f build/frontend.Dockerfile \
-t $ACR_NAME.azurecr.io/dapr-store/frontend-host

docker push $ACR_NAME.azurecr.io/dapr-store/products
docker push $ACR_NAME.azurecr.io/dapr-store/orders
docker push $ACR_NAME.azurecr.io/dapr-store/frontend-host