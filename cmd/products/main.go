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

	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/env"

	"database/sql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
	_ "github.com/mattn/go-sqlite3"
)

// API type is a wrap of he common API with local version
type API struct {
	*api.APIBase
}

var (
	healthy     = true               // Simple health flag
	version     = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName = "products"
	defaultPort = 9002
	daprPort    int
	db          *sql.DB
)

//
// Main entry point, will start HTTP service
//
func main() {
	log.SetOutput(os.Stdout) // Personal preference on log output
	log.Printf("### Dapr Store: %v v%v starting...", serviceName, version)

	// Port to listen on, change the default as you see fit
	serverPort := env.GetEnvInt("PORT", defaultPort)

	daprPort = env.GetEnvInt("DAPR_HTTP_PORT", 0)
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
		&api.APIBase{
			ServiceName: serviceName,
			Healthy:     healthy,
			Version:     version,
			BuildInfo:   buildInfo,
		}}

	api.AddCommonRoutes(router)
	api.addRoutes(router)

	var err error
	db, err = sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Panicf("### Failed to open database! %+v\n", err)
	}
	defer db.Close()

	// Start server
	log.Printf("### Server listening on %v\n", serverPort)
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}
	err = srv.ListenAndServe()
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
