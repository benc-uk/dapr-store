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
The major components and microservices that make up the Dapr Store system are described here
>Note. The term "component" here, is used in the general english sense, rather than referring to a [*Dapr component*](https://github.com/dapr/docs/tree/master/concepts#components) 

## Shared Go packages
Shared Go code lives in the `pkg/` directory, which is used by all the services. These fall into the following packages:
- `pkg/api` - A base API extended by all services, provides health & status endpoints.
- `pkg/auth` - Server side token validation of JWT using JWK.
- `pkg/env` - Very simple `os.LookupEnv` wrapper with fallback defaults.
- `pkg/models` - Types and data structs used by the services.
- `pkg/problem` - Standarized REST error messages using [RFC 7807 Problem Details](https://tools.ietf.org/html/rfc7807).
- `pkg/state` - A Dapr state API helper for getting / setting state with error handling.

## Orders service
This service provides order processing to the Dapr Store.  
It is written in Go, source is in `cmd/orders` and it exposes the following API routes:
```
/get/{id}                GET a single order by orderID
/getForUser/{username}   GET all orders for a given username
```
See `pkg/models` for details of the **Order** struct.

The service provides some faked order processing activity so that orders are moved through a number of statuses, simulating some back-office systems or inventory management. Orders are initially set to `OrderReceived` status, then after 30 seconds moved to `OrderProcessing`, then after 2 minutes moved to `OrderComplete`

### Orders - Dapr Interaction
- **Pub/Sub.** Subscribes to the `orders-queue` topic to receive new orders from the *cart* service
- **State.** Stores and retrieves **Order** entities from the state service, keyed on OrderID. Also lists of orders per user, held as an array of OrderIDs and keyed on username
- **SendGrid.** To be added


## Users service
This provides a simple user profile service to the Dapr Store. Only registered users can use the store to place orders etc.  
It is written in Go, source is in `cmd/users` and it exposes the following API routes:
```
/register                 POST a new user to register them
/get/{username}           GET the user profile for given username
/isregistered/{username}  GET the registration status for a given username
```
See `pkg/models` for details of the **User** struct.

The service provides some faked order processing activity so that orders are moved through a number of statuses, simulating some back-office systems or inventory management. Orders are initially set to `OrderReceived` status, then after 30 seconds moved to `OrderProcessing`, then after 2 minutes moved to `OrderComplete`

### Users - Dapr Interaction
- **State.** Stores and retrieves **User** entities from the state service, keyed on username.


## Products service
This is the product catalog service for Dapr Store.  
It is written in Go, source is in `cmd/products` and it exposes the following API routes:
```
/get/{id}        GET a single product with given id
/catalog         GET all products in the catalog, returns an array of products
/offers          GET all products that are on offer, returns an array of products
/search/{query}  GET search the product database, returns an array of products
```
See `pkg/models` for details of the **Product** struct.

The products data is held in a SQLite database, this decision was taken due to the lack of support for queries and filtering with the Dapr state API. The source data to populate the DB is in `etc/products.csv` and the database can be created with the `scripts/create-products-db.sh`. The database file (sqlite.db) is currently stored inside the products container, effectively making catalogue baked in at build time. This could be changed/improved at a later date

### Products - Dapr Interaction
None


## Cart service
This provides a cart service to the Dapr Store. The currently implementation is very minimal.  
It is written in Go, source is in `cmd/cart` and it exposes the following API routes:
```
/submit    POST a new order to be submitted
```
See `pkg/models` for details of the **Order** struct.

The service currently is little more than gateway API for publishing new orders onto the `orders-queue`. The cart is not persisted server side (it's held in the client in local storage), this leaves room for future improvement.

### Users - Dapr Interaction
- **Pub/Sub.** Pushes **Order** entities to the `orders-queue` topic to be collected by the *orders* service

## Frontend 
This is the frontend accessed by users of store and visitors to the site. It is a single-page application (SPA) as such it runs entirely client side in the browser. It was created using the [Vue CLI](https://cli.vuejs.org/) and written in Vue.js

It follows the standard SPA pattern of being served via static hosting (the 'frontend host' described below) and all data is fetched via a REST API endpoint. Note. [Vue Router](https://router.vuejs.org/) is used to provide client side routing, as such it needs to be served from a host that is configured to support it.

The default API endpoint is `/` and it makes calls to the Dapr invoke API, namely `/v1.0/invoke/{service}` via the API gateway

## Frontend host
A very basic / minimal static server using gorilla/mux. See https://github.com/gorilla/mux#static-files. It simply serves up the static bundled files output from the build process of the frontend, it expects to find these files in `./dist` directory but this is configurable 

## API gateway



# Running in Kubernetes  - Quick guide 

# Running locally - Quick guide 
This is a (very) basic guide to running the Dapr Store locally. Only instructions for WSL 2/Linux/MacOS are provided. It's advised to only do this if you wish to develop or debug the project.

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
The services support the following environmental variables. All settings are optional.

- `PORT` - Port the server will listen on. See defaults below.
- `AUTH_CLIENT_ID` - Used to enable [security and authentication](/docs/security.md). Default is *blank*, which is no security
- `DAPR_STORE_NAME` - Name of the Dapr state component to use. Default is `statestore`
- `DAPR_ORDERS_TOPIC` - Name of the Dapr pub/sub topic to use for orders. Default is `orders-queue"`
- `STATIC_DIR` - (Frontend host only) The path to serve static content from, i.e. the bundled Vue.js SPA output. Default is `./dist`

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