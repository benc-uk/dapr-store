// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for products
// ----------------------------------------------------------------------------

package main

import (
	"log"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/cmd/products/impl"
	"github.com/benc-uk/dapr-store/cmd/products/spec"

	"github.com/benc-uk/go-rest-api/pkg/api"
	"github.com/benc-uk/go-rest-api/pkg/env"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
	_ "github.com/mattn/go-sqlite3"
)

// API type is a wrap of the common base API with local implementation
type API struct {
	*api.Base
	service spec.ProductService
}

var (
	healthy     = true               // Simple health flag
	version     = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName = "products"
	defaultPort = 9002
)

// Main entry point, will start HTTP service
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output

	// Port to listen on, change the default as you see fit
	serverPort := env.GetEnvInt("PORT", defaultPort)

	// Use chi for routing
	router := chi.NewRouter()

	dbFilePath := "./sqlite.db"
	if len(os.Args) > 1 {
		dbFilePath = os.Args[1]
	}

	// Wrapper API with anonymous inner new Base API
	api := API{
		api.NewBase(serviceName, version, buildInfo, healthy),
		impl.NewService(serviceName, dbFilePath),
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

	// Add application routes for this service
	api.addRoutes(router)

	// Finally start the server
	api.StartServer(serverPort, router, 5*time.Second)
}
