# Dapr Store

Dapr Store is a sample/reference application showcasing the use of [Dapr](https://dapr.io/) to build microservices based applications. It is a simple online store with all the core components that make up such system, e.g. a frontend for users, authentication, product catalog, and order processing etc.

[Dapr](https://dapr.io/) is an _"event-driven, portable runtime for building microservices on cloud and edge"_. The intention of this project was to show off many of the capabilities and features of Dapr, but in the context of a real working application. This has influenced the architecture and design decisions, balancing between realism and a simple _"demo-ware"_ showcase.

The backend microservices are written in Go (however it's worth nothing that Dapr is language independent), and the frontend is a single-page application (SPA) written in [Vue.js](https://vuejs.org/). All APIs are REST & HTTP based

This repo is a monorepo, containing the source for several discreet but closely linked codebases, one for each component of the project, as described below.  
The ["Go Standard Project Layout"](https://github.com/golang-standards/project-layout) has been used.

:warning: The project is still in the experimental stage, with a high rate of change - expect breaking changes

# Architecture

The following diagram shows all the components of the application and main interactions. It also highlights which Dapr API/feature (aka Dapr building block) is used and where.
![architecture diagram](./docs/img/design.png)

## Dapr Interfaces & Building Blocks

The application uses the following [Dapr Building Blocks](https://docs.dapr.io/developing-applications/building-blocks/) and APIs

- **Service Invocation** — The API gateway calls the four main microservices using HTTP calls to [Dapr service invocation](https://docs.dapr.io/developing-applications/building-blocks/service-invocation/service-invocation-overview/). This provides retries, mTLS and service discovery.
- **State** — State is held for users and orders using the [Dapr state management API](https://docs.dapr.io/developing-applications/building-blocks/state-management/state-management-overview/). The state provider used is Redis, however any other provider could be plugged in without any application code changes.
- **Pub/Sub** — The submission of new orders through the cart service, is decoupled from the order processing via pub/sub messaging and the [Dapr pub/sub messaging API](https://docs.dapr.io/developing-applications/building-blocks/pubsub/pubsub-overview/). New orders are placed on a topic as messages, to be collected by the orders service. This allows the orders service to independently scale and separates our reads & writes
- **Output Bindings** — To communicate with downstream & 3rd party systems, the [Dapr Bindings API](https://docs.dapr.io/developing-applications/building-blocks/bindings/bindings-overview/) is used. This allows the store to carry out tasks such as saving order details into external storage (e.g. Azure Blob) and notify uses with emails via SendGrid
- **Middleware** — Dapr supports a range of HTTP middleware, for this project traffic rate limiting can enabled on any of the APIs with a single Kubernetes annotation

# Project Status

![](https://img.shields.io/github/last-commit/benc-uk/dapr-store) ![](https://img.shields.io/github/release-date/benc-uk/dapr-store) ![](https://img.shields.io/github/v/release/benc-uk/dapr-store) ![](https://img.shields.io/github/commit-activity/m/benc-uk/dapr-store)

Deployed instance: https://daprstore.kube.benco.io/  
[![](https://img.shields.io/website?label=Hosted%3A%20Kubernetes&up_message=online&url=https%3A%2F%2Fdaprstore.kube.benco.io%2F)](https://daprstore.kube.benco.io/)

# Application Elements & Services

The main elements and microservices that make up the Dapr Store system are described here

## Shared Go packages

Shared Go code lives in the `pkg/` directory, which is used by all the services, and consists of the following packages:

- `pkg/api` - A base API extended by all services, provides health & status endpoints.
- `pkg/apitests` - A simple helper for running sets of router/API based tests.
- `pkg/auth` - Server side token validation of JWT using JWK.
- `pkg/dapr` - A Dapr helper & wrapper library for state, pub/sub and output bindings
- `pkg/env` - Very simple `os.LookupEnv` wrapper with fallback defaults.
- `pkg/problem` - Standarized REST error messages using [RFC 7807 Problem Details](https://tools.ietf.org/html/rfc7807).

## Service Code

Each Go microservice (in `cmd/`) follows a very similar layout to the code base (the exception being `frontend-host` which has no business logic)

Primary runtime code:

- `main.go` - Starts HTTP server, creates service implementation + main entry point
- `routes.go` - All controllers for routes exposed by the service's API
- `spec/spec.go` - Specification of domain entity (e.g. User) and interface to support it
- `impl/impl.go` - Concrete implementation of the above spec, backed by either Dapr or other dependency (e.g. a database)

For testing:

- `*_test.go` - Main service tests
- `mock/mock.go` - Mock implementation of the domain spec, with no dependencies

## 💰 Orders service

This service provides order processing to the Dapr Store.  
It is written in Go, source is in `cmd/orders` and it exposes the following API routes:

```
/get/{id}                GET a single order by orderID
/getForUser/{username}   GET all orders for a given username
```

See `cmd/orders/spec` for details of the **Order** entity.

The service provides some fake order processing activity so that orders are moved through a number of statuses, simulating some back-office systems or inventory management. Orders are initially set to `OrderReceived` status, then after 30 seconds moved to `OrderProcessing`, then after 2 minutes moved to `OrderComplete`

### Orders - Dapr Interaction

- **Pub/Sub.** Subscribes to the `orders-queue` topic to receive new orders from the cart service
- **State.** Stores and retrieves **Order** entities from the state service, keyed on OrderID. Also lists of orders per user, held as an array of OrderIDs and keyed on username
- **Bindings.** All output bindings are optional, the service operates without these present
  - **Azure Blob.** For saving "order reports" as text files into Azure Blob storage
  - **SendGrid.** For sending emails to users via [SendGrid](https://sendgrid.com/)

## 👦 Users service

This provides a simple user profile service to the Dapr Store. Only registered users can use the store to place orders etc.  
It is written in Go, source is in `cmd/users` and it exposes the following API routes:

```
/register                 POST a new user to register them
/get/{username}           GET the user profile for given username
/isregistered/{username}  GET the registration status for a given username
```

See `cmd/users/spec` for details of the **User** entity.

The service is notable as it consists of a mixture of secured API routes and one anonymous/open API `/isregistered`

### Users - Dapr Interaction

- **State.** Stores and retrieves **User** entities from the state service, keyed on username.

## 📑 Products service

This is the product catalog service for Dapr Store.  
It is written in Go, source is in `cmd/products` and it exposes the following API routes:

```
/get/{id}        GET a single product with given id
/catalog         GET all products in the catalog, returns an array of products
/offers          GET all products that are on offer, returns an array of products
/search/{query}  GET search the product database, returns an array of products
```

See `cmd/products/spec` for details of the **Product** entity.

The products data is held in a SQLite database, this decision was taken due to the lack of support for queries and filtering with the Dapr state API. The source data to populate the DB is in `etc/products.csv` and the database can be created with the `scripts/create-products-db.sh`. The database file (sqlite.db) is currently stored inside the products container, effectively making catalogue baked in at build time. This could be changed/improved at a later date

### Products - Dapr Interaction

None directly, but is called via service invocation from other services, the API gateway & the cart service.

## 🛒 Cart service

This provides a cart service to the Dapr Store. The currently implementation is a MVP.  
It is written in Go, source is in `cmd/cart` and it exposes the following API routes:

```
/setProduct/{username}/{productId}/{count}    PUT a number of products in the cart of given user
/get/{username}                               GET cart for user
/submit                                       POST submit a cart, and turn it into an 'Order'
/clear/{username}                             PUT clear a user's cart

```

See `pkg/models` for details of the **Order** struct.

The service is responsible for maintaining shopping carts for each user and persisting them. Submitting a cart will validate the contents and turn it into a order, which is sent to the Orders service for processing

### Cart - Dapr Interaction

- **Pub/Sub.** The cart pushes **Order** entities to the `orders-queue` topic to be collected by the orders service
- **State.** Stores and retrieves **Cart** entities from the state service, keyed on username.
- **Service Invocation.** Cross service call to products API to lookup and check products in the cart

## 💻 Frontend

This is the frontend accessed by users of store and visitors to the site. It is a single-page application (SPA) as such it runs entirely client side in the browser. It was created using the [Vue CLI](https://cli.vuejs.org/) and written in Vue.js

It follows the standard SPA pattern of being served via static hosting (the 'frontend host' described below) and all data is fetched via a REST API endpoint. Note. [Vue Router](https://router.vuejs.org/) is used to provide client side routing, as such it needs to be served from a host that is configured to support it.

The default API endpoint is `/` and it makes calls to the Dapr invoke API, namely `/v1.0/invoke/{service}` this is routed via the _API gateway_ to the various services.

## 📡 Frontend host

A very standard static content server using gorilla/mux. See https://github.com/gorilla/mux#static-files. It simply serves up the static bundled files output from the build process of the frontend, it expects to find these files in `./dist` directory but this is configurable

In addition it exposes a simple `/config` endpoint, this is to allow dynamic configuration of the frontend. It passes two env vars `AUTH_CLIENT_ID` and `API_ENDPOINT` from the frontend host to the frontend Vue SPA as a JSON response, which are fetched & read as the app is loaded in the browser.

## 🌍 API gateway

This component is critical but consists of no code. It's a NGINX reverse proxy configured to do two things:

- Forward specific calls to the relevant services via Dapr
- Direct requests to the _frontend host_

> Note. This is not to be confused with Azure API Management, Azure App Gateway or AWS API Gateway 😀

This is done with path based routing, it aggregates the various APIs and frontend SPA into a single endpoint or host, making configuration much easier (and the reason the API endpoint for the SPA can simply be `/`)

NGINX is run with the Dapr sidecar alongside it, so that it can proxy requests to the `/v1.0/invoke` Dapr API, via this the downstream services are invoked, through Dapr.

### Within Kubernetes

Inspired by [this blog post](https://carlos.mendible.com/2020/04/05/kubernetes-nginx-ingress-controller-with-dapr/) it is deployed in Kubernetes as a "Daprized NGINX ingress controller". See `deploy/ingress` for details on how this is done.

### Locally

To provide a like for like experience with Kubernetes, and the single aggregated endpoint - the same model is used locally. NGINX is run as a Docker container exposed to the host network, and NGINX configuration applied allow to route traffic. This is then run via the dapr CLI (i.e `dapr run`) so that the daprd sidecar process is available to it.

See `scripts/local-gateway` for details on how this is done, the `scripts/local-gateway/run.sh` script starts the gateway which will run on port 9000

# Running in Kubernetes - Quick guide

See [deploy/readme.md](deploy/readme.md)

# Running Locally - Quick guide

This is a (very) basic guide to running Dapr Store locally. Only instructions for WSL 2/Linux/MacOS are provided. It's advised to only do this if you wish to develop or debug the project.

### Prereqs

- Docker
- Go v1.15+
- Node.js v12+

### Setup

Install and initialize Dapr

```
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
dapr init
```

### Clone repo

```bash
git clone https://github.com/benc-uk/dapr-store/
```

### Run all services

Run everything. Run from the project root (e.g. `dapr-store` directory), this will run all the services, the API gateway and the Vue frontend

```bash
make run
```

Access the store from http://localhost:9000/

**💣 Gotcha!** The Vue frontend will start and display a message "App running at" saying it is running on port 8000, **do not access the frontend directly on this port, it will not function!**, always go via the gateway running on port 9000

# Working Locally

A makefile is provided to assist working with the project and building/running it, the main targets are:

```text
help                 💬 This help message :)
lint                 🔎 Lint & format, check to be run in CI, sets exit code on error
lint-fix             📝 Lint & format, fixes errors and modifies code
test                 🎯 Unit tests for services and snapshot tests for SPA frontend
test-reports         📜 Unit tests with coverage and test reports
test-snapshot        📷 Update snapshots for frontend tests
image-all            📦 Build all container images
push-all             📤 Push all images to registry
bundle               💻 Build and bundle the frontend Vue SPA
clean                🧹 Clean the project, remove modules, binaries and outputs
run                  🚀 Start & run everything locally
stop                 ⛔ Stop & kill everything started locally from `make run`
```

# CI / CD

A working set of CI and CD release GitHub Actions workflows are provided `.github/workflows/`, automated builds are run in GitHub hosted runners

### [GitHub Actions](https://github.com/benc-uk/dapr-store/actions)

[![](https://img.shields.io/github/workflow/status/benc-uk/dapr-store/CI%20Build%20App/master?label=CI+Build+App)](https://github.com/benc-uk/dapr-store/actions/workflows/ci-build.yml)

[![](https://img.shields.io/github/workflow/status/benc-uk/dapr-store/Deploy%20To%20Kubernetes/master?label=Deploy+to+Kubernetes)](https://github.com/benc-uk/dapr-store/actions/workflows/deploy-k8s.yaml)

# Security, Identity & Authentication

The default mode of operation for the Dapr Store is in "demo mode" where there is no identity provided configured, and no security on the APIs. This makes it simple to run and allows us to focus on the Dapr aspects of the project. In this mode a demo/dummy user account can be used to sign-in and place orders in the store.

Optionally Dapr store can be configured utilise the [Microsoft identity platform](https://docs.microsoft.com/en-us/azure/active-directory/develop/) (aka Azure Active Directory v2) as an identity provider, to enable real user sign-in, and securing of the APIs.

See the [security, identity & authentication docs](./docs/auth-identity.md) for more details on setting this up.

# Configuration

## Environmental Variables

The services support the following environmental variables. All settings are optional.

- `PORT` - Port the server will listen on. See defaults below.
- `AUTH_CLIENT_ID` - Used to enable integration with Azure AD for identity and authentication. Default is _blank_, which runs the service with no identity backend. See the [security, identity & authentication docs](./docs/auth-identity.md) for more details.
- `DAPR_STORE_NAME` - Name of the Dapr state component to use. Default is `statestore`
- `DAPR_ORDERS_TOPIC` - Name of the Dapr pub/sub topic to use for orders. Default is `orders-queue`
- `DAPR_PUBSUB_NAME` - Name of the Dapr pub/sub component to use for orders. Default is `pubsub`

Frontend host config:

- `STATIC_DIR` - The path to serve static content from, i.e. the bundled Vue.js SPA output. Default is `./dist`
- `API_ENDPOINT` - To point the frontend at a different endpoint. It's very unlikely you'll ever need to set this. Default is `/`

## Default ports

- 9000 - NGINX API gateway (reverse proxy)
- 9001 - Cart service
- 9002 - Products service
- 9003 - Users service
- 9004 - Order processing service
- 8000 - Frontend host

## Dapr Components

See the [components documentation](components/) for full details of the Dapr components used by the application and how to configure them.

# Roadmap & known issues

See [project plan on GitHub](https://github.com/benc-uk/dapr-store/projects/1)

# Concepts and Terms

Clarity of terminology is sometimes important, here's a small glossary

- **Building Block** - [Specific Dapr term](https://github.com/dapr/docs/tree/master/concepts#building-blocks). A _building block_ is an API level feature of Dapr, such as 'state management' or 'pub/sub' or 'secrets'
- **Component** - Component is another Dapr specific term. A _component_ is a plugin that provides implementation functionality to building blocks. As component is a generic & commonly used word, the term "Dapr component" will be used where ambiguity is possible
- **Service** - The microservices, written in Go and exposing REST API, either invoked through Dapr and.or using the Dapr API for things such as state.
- **API Gateway** - NGINX reverse proxy sitting in front of the services. This is not to be confused with Azure API Management, Azure App Gateway or AWS API Gateway
- **State** - Dapr state API, backed with a Dapr component state provider, e.g. Redis
- **Entity** - A data object, typically a JSON representation of one of the structs in `pkg/models`, can be client or server side
