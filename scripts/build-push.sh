#!/bin/bash
set -e

if [[ -z $1 ]]; then
  echo "Error! Please provide Docker registry and/or repo, e.g. 'mydockeruser' or 'myreg.azurecr.io/daprstore'"
  exit 1
fi

echo -e "\n\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo -e "\e[34m»»» 🚀 \e[32mBuilding \e[33m'cart'\e[32m service ..."
echo -e "\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»\u001b[0m"
docker build . -f build/service.Dockerfile \
--build-arg serviceName=cart \
--build-arg servicePort=9001 \
-t $1/cart

echo -e "\n\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo -e "\e[34m»»» 🚀 \e[32mBuilding \e[33m'products'\e[32m service ..."
echo -e "\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»\u001b[0m"
docker build . -f build/service.Dockerfile \
--build-arg serviceName=products \
--build-arg servicePort=9002 \
-t $1/products

echo -e "\n\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo -e "\e[34m»»» 🚀 \e[32mBuilding \e[33m'users'\e[32m service ..."
echo -e "\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»\u001b[0m"
docker build . -f build/service.Dockerfile \
--build-arg serviceName=users \
--build-arg servicePort=9003 \
-t $1/users

echo -e "\n\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo -e "\e[34m»»» 🚀 \e[32mBuilding \e[33m'orders'\e[32m service ..."
echo -e "\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»\u001b[0m"
docker build . -f build/service.Dockerfile \
--build-arg serviceName=orders \
--build-arg servicePort=9004 \
-t $1/orders

echo -e "\n\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo -e "\e[34m»»» 🚀 \e[32mBuilding \e[33m'frontend-host'\e[32m service ..."
echo -e "\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»\u001b[0m"
docker build . -f build/frontend.Dockerfile \
-t $1/frontend-host

echo -e "\n\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»"
echo -e "\e[34m»»» 🪐 \e[32mPushing images to \e[33m'$1'\e[32m ..."
echo -e "\e[34m»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»»\u001b[0m"
docker push $1/cart
docker push $1/products
docker push $1/users
docker push $1/orders
docker push $1/frontend-host