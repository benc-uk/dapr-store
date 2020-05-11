#!/bin/bash

#
# This script runs NGINX as a reverse proxy as a docker container
# It expects to find DAPR_HTTP_PORT variable, so always run via `dapr run`
#

echo "### API gateway nginx wrapper magic doohicky script"
echo "### Dapr port is $DAPR_HTTP_PORT"
echo "### Starting nginx in docker"

sed -i "s/localhost:[[:digit:]]\+; #dapr/localhost:$DAPR_HTTP_PORT; #dapr/g" dapr.conf

# Disabled this nonsense for now, it was an attempt to get it working in a devcontainer
# LOCAL_GW_DIR=$PWD
# if [ ! -z "$PROJECT_PATH" ]; then
#   LOCAL_GW_DIR="$PROJECT_PATH/scripts/local-gateway"
# fi

docker run --rm --network host -p 9000:9000 --mount type=bind,source=$PWD,target=/etc/nginx/conf.d/ --name api-gateway nginx:alpine