# Dapr Store

This is a sample application showcasing the use of [Dapr](https://dapr.io/) to build microservices based systems

It is written in a mixture of Go and Vue.js

![architecture diagram](./docs/img/design.png)


## Run Locally - Quick Guide (WSL 2/Linux/MacOS)

### Prereqs
- Docker
- Go 1.14+
- Realize, install with: `go get github.com/oxequa/realize`

### Setup
Install and initialize Dapr
```
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
sudo dapr init
```

### Run All Services
Run everything locally.
From root of project (e.g. `dapr-store` directory)
```
find . -name '*.sh' -print0 |xargs -0 chmod +x
./bin/start-local.sh
```

To stop Dapr instances and other processes, run the `start-local.sh` script:
```
./bin/stop-local.sh
```
# Reference Information

## Config - Environmental Variables
- `PORT`
- `AUTH_CLIENT_ID`
- `DAPR_STORE_NAME`
- `DAPR_ORDERS_TOPIC`
- `STATIC_DIR` (Front host only)

## Local & default ports
- 9000 - NGINX API gateway (reverse proxy)
- 9001 - Cart service
- 9002 - Products service
- 9003 - Users service
- 9004 - Order processing service
- 8000 - Frontend host
