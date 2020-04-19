# Dapr Store
Dapr Store is a sample/reference application showcasing the use of [Dapr](https://dapr.io/) to build microservices based systems. It is a simple online store with all the core elements that make up such system, e.g. a frontend for users, authentication, product catalog, and order processing

[Dapr](https://dapr.io/) is an *"event-driven, portable runtime for building microservices on cloud and edge"*. The intention of this project was to show off many of the capabilities and features of Dapr in the context of a real working application. This has influenced the architecture and design decisions, balancing between realism and a simple demo-ware showcase.

The backend microservices are written in Go (however it's worth nothing that Dapr is language independent), and the frontend is a single-page application (SPA) written in [Vue.js](https://vuejs.org/). All APIs are REST & HTTP based

This repo is a monorepo, containing the source for several discreet but closely linked codebases, one for each component of the project, as described below.  
The ["Go Standard Project Layout"](https://github.com/golang-standards/project-layout) has been used. 


# Architecture
The following diagram shows all the components of the application and main interactions. It also highlights which Dapr API/feature (aka Dapr building block) is used and where.
![architecture diagram](./docs/img/design.png)


# Components

## Cart service
## Orders service
## Users service
## Products service
## Frontend 
This is the frontend accessed by users of store and visitors to the site. It is a single-page application (SPA) as such it runs entirely client side in the browser.

It follows the standard SPA pattern of being served via static hosting (the 'frontend host' described below) and all data is fetched live, via a REST API

Vue router 

## Frontend host
## API gateway

>Note. The term "component" here, is used in the general english sense, rather than referring to a [*Dapr component*](https://github.com/dapr/docs/tree/master/concepts#components) 


# Running locally - Quick guide 
This is a (very) basic guide to running the Dapr Store locally. Only instructions for WSL 2/Linux/MacOS are provided

### Prereqs
- Docker
- Go 1.14+
- Realize. install with: `go get github.com/oxequa/realize`. Realize is a task runner for Go, with live reloading

### Setup
Install and initialize Dapr
```
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
sudo dapr init
```
### Clone repo
```bash
git clone https://github.com/benc-uk/dapr-store/
# Make all scripts executable (only needed after first clone)
cd dapr-store
find . -name '*.sh' -print0 |xargs -0 chmod +x
```

### Run all services
Run everything. Run from the project root (e.g. `dapr-store` directory)
```bash
./bin/start-local.sh
```

To stop Dapr instances and other processes, run the `stop-local.sh` script:
```bash
./bin/stop-local.sh
```


# Reference Information
NOT FINISHED ☢

## Config - Environmental Variables
- `PORT`
- `AUTH_CLIENT_ID`
- `DAPR_STORE_NAME`
- `DAPR_ORDERS_TOPIC`
- `STATIC_DIR` (Front host only)

## Default ports
- 9000 - NGINX API gateway (reverse proxy)
- 9001 - Cart service
- 9002 - Products service
- 9003 - Users service
- 9004 - Order processing service
- 8000 - Frontend host


# Roadmap & known issues
See [project plan on GitHub](https://github.com/benc-uk/dapr-store/projects/1)


# Concepts and Terms
Clarity of terminology is important

- **Building Block** - [Specific Dapr term](https://github.com/dapr/docs/tree/master/concepts#building-blocks). A *building bloc*k is an API level feature of Dapr, such as 'state mangement' or 'pub/sub' or 'secrets'
- **Component** - Component is another Dapr specific term. A *component* is a plugin that provides implementation functionality to building blocks. As component is a generic & commonly used word, the term "Dapr component" will be used where ambiguity is possible
- **Service** - The 
- API Gateway 
- State