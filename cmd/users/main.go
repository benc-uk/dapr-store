// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for users
// ----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/dapr"
	"github.com/benc-uk/dapr-store/pkg/env"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// API type is a wrap of he common API with local version
type API struct {
	*api.APIBase
}

var (
	healthy     = true               // Simple health flag
	version     = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName = "users"
	defaultPort = 9003
	daprHelper  *dapr.Helper
)

//
// Main entry point, will start HTTP service
//
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output
	log.Printf("### Dapr Store: %v v%v starting...", serviceName, version)

	// Port to listen on, change the default as you see fit
	serverPort := env.GetEnvInt("PORT", defaultPort)

	// Bootstrap standard helper, checks env vars for default settings etc
	daprHelper = dapr.BootstrapHelper(serviceName)
	if daprHelper == nil {
		os.Exit(1)
	}

	// Use gorilla/mux for routing
	router := mux.NewRouter()

	// Add middleware for logging and CORS
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	// Wrapper type with anonymous inner field
	api := API{
		&api.APIBase{
			ServiceName: serviceName,
			Healthy:     healthy,
			Version:     version,
			BuildInfo:   buildInfo,
		}}

	api.AddCommonRoutes(router)
	api.addRoutes(router)

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
	return handlers.LoggingHandler(os.Stdout, handler)
}
