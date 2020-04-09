// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service
// ----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benc-uk/dapr-store/common"
	"github.com/benc-uk/go-starter/pkg/envhelper"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// API type is a wrap of he common API with local version
type API struct {
	*common.API
}

var (
	healthy        = true               // Simple health flag
	version        = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo      = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName    = "products"
	daprPort       int
	daprStateStore string
)

//
// Main entry point, will start HTTP service
//
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output
	log.Printf("### Dapr Store: %v v%v starting...", serviceName, version)

	// Port to listen on, change the default as you see fit
	serverPort := envhelper.GetEnvInt("PORT", 9002)
	daprStateStore = envhelper.GetEnvString("DAPR_STORE_NAME", "statestore")

	daprPort = envhelper.GetEnvInt("DAPR_HTTP_PORT", 0)
	if daprPort != 0 {
		log.Printf("### Dapr sidecar detected on port %v", daprPort)
	} else {
		log.Printf("### Dapr not detected (no DAPR_HTTP_PORT available), this is bad")
		log.Printf("### Exiting...")
		os.Exit(1)
	}

	// Use gorilla/mux for routing
	router := mux.NewRouter()

	// Add middleware for logging and CORS
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	// Wrapper type with anonymous inner field
	api := API{
		&common.API{
			ServiceName: serviceName,
			Healthy:     healthy,
			Version:     version,
			BuildInfo:   buildInfo,
		}}

	api.AddCommonRoutes(router)
	api.addRoutes(router)

	// Start server
	log.Printf("### Dapr state store: %v\n", daprStateStore)
	log.Printf("### Server listening on %v\n", serverPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
	if err != nil {
		panic(err.Error())
	}
}

//
// Change CORS settings here
//
func corsMiddleware(handler http.Handler) http.Handler {
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	return handlers.CORS(corsOrigins, corsMethods)(handler)
}

//
// Change request logging here
//
func loggingMiddleware(handler http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, handler)
}
