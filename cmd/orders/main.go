// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for orders
// ----------------------------------------------------------------------------

package main

import (
	"log"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/impl"
	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/go-rest-api/pkg/api"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/dapr/pubsub"
	"github.com/benc-uk/go-rest-api/pkg/env"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// API type is a wrap of the common base API with local implementation
type API struct {
	*api.Base
	service spec.OrderService
}

var (
	healthy     = true               // Simple health flag
	version     = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName = "orders"
	defaultPort = 9004
)

// Main entry point, will start HTTP service
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output

	// Port to listen on, change the default as you see fit
	serverPort := env.GetEnvInt("PORT", defaultPort)

	// Needed for pub sub
	pubSubName := env.GetEnvString("DAPR_PUBSUB_NAME", "pubsub")
	topicName := env.GetEnvString("DAPR_ORDERS_TOPIC", "orders-queue")

	// Use chi for routing
	router := chi.NewRouter()

	svc := impl.NewService(serviceName)
	// Wrapper API with anonymous inner new Base API
	api := API{
		api.NewBase(serviceName, version, buildInfo, healthy),
		svc,
	}

	// Enabling of auth is optional, set via AUTH_CLIENT_ID env var
	var validator auth.Validator

	if clientID := env.GetEnvString("AUTH_CLIENT_ID", ""); clientID == "" {
		log.Println("### 🚨 No AUTH_CLIENT_ID set, API auth will be disabled")

		validator = auth.NewPassthroughValidator()
	} else {
		log.Println("### 🔐 Auth enabled, API will be protected with JWT validation")

		validator = auth.NewJWTValidator(clientID, "https://login.microsoftonline.com/common/discovery/v2.0/keys", "store-api")
	}

	// Some basic middleware
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// Some custom middleware for CORS
	router.Use(api.SimpleCORSMiddleware)

	// Add Prometheus metrics endpoint, must be before the other routes
	api.AddMetricsEndpoint(router, "metrics")

	// Add root, health & status middleware
	api.AddHealthEndpoint(router, "health")
	api.AddStatusEndpoint(router, "status")
	api.AddOKEndpoint(router, "")

	// Special Dapr endpoints added to the router to support pub/sub
	pubsub.Subscribe(pubSubName, []string{topicName}, router)
	pubsub.AddTopicHandler(topicName, router, svc.PubSubOrderReceiver)

	// Add application routes for this service
	api.addRoutes(router, validator)

	// Finally start the server
	api.StartServer(serverPort, router, 5*time.Second)
}
