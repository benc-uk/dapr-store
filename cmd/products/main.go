// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for products
// ----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/cmd/products/impl"
	"github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/env"

	"github.com/gorilla/mux"
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

//
// Main entry point, will start HTTP service
//
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output
	log.Printf("### Dapr Store: %v v%v starting...", serviceName, version)

	// Port to listen on, change the default as you see fit
	serverPort := env.GetEnvInt("PORT", defaultPort)

	// Use gorilla/mux for routing
	router := mux.NewRouter()

	dbFilePath := "./sqlite.db"
	if len(os.Args) > 1 {
		dbFilePath = os.Args[1]
	}

	// Wrapper API with anonymous inner new Base API
	api := API{
		api.NewBase(serviceName, version, buildInfo, healthy, router),
		impl.NewService(serviceName, dbFilePath),
	}

	// Add routes for this service
	api.addRoutes(router)

	// Add middleware routes, these are all optional
	api.AddStatus(router)  // Add status and information endpoint
	api.AddLogging(router) // Add request logging
	api.AddHealth(router)  // Add health endpoint
	api.AddMetrics(router) // Expose metrics, in prometheus format
	api.AddRoot(router)    // Respond to root request with a simple 200 OK

	// Start server
	log.Printf("### Server listening on %v\n", serverPort)
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
