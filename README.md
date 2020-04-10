# Dapr Store

This is a sample application showcasing the use of [Dapr](https://dapr.io/) to build microservices based systems

It is written in a mixture of Go and Vue.js

![architecture diagram](./docs/img/design.png)


## Run Locally - Quick Guide (WSL/Linux)

### Prereqs
- Docker
- Go 1.14+
- Realize, install with: `go get github.com/oxequa/realize`

Install and initialize Dapr
```
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
sudo dapr init
```

Run everything locally.
From root of project (e.g. `dapr-store` directory)
```
find . -name '*.sh' -print0 |xargs -0 chmod +x
./bin/start-local.sh
```

To stop Dapr instances and other processes, press Ctrl+C to exit the `start-local.sh` script, then run:
```
./bin/stop-local.sh
```
# Reference Information

## Local ports
9000 - NGINX API gateway (reverse proxy)
9001 - Orders service
9001 - Products service
8000 - Frontend host
