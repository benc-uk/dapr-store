#!/bin/bash
scriptDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

#
# This script runs NGINX as a reverse proxy as a docker container
# It expects to find DAPR_HTTP_PORT variable, so always run via `dapr run`
#

echo "🠶🠶🠶 🌐 API gateway nginx wrapper magic script thing"
echo "🠶🠶🠶 🎩 Dapr port is $DAPR_HTTP_PORT"
echo "🠶🠶🠶 📦 Starting nginx in docker"

# This voodoo dynamically sets & overwrites the port for the Dapr process in the nginx config file
sed -i "s/localhost:[[:digit:]]\+; #dapr/localhost:$DAPR_HTTP_PORT; #dapr/g" $scriptDir/dapr.conf

# The starts nginx in Docker
docker run --rm --network host --mount type=bind,source=$scriptDir,target=/etc/nginx/conf.d/ --name api-gateway nginx:alpine