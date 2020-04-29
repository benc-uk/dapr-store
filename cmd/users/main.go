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

	"github.com/benc-uk/dapr-store/cmd/users/impl"
	"github.com/benc-uk/dapr-store/cmd/users/spec"
	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/env"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// API type is a wrap of the common base API with local implementation
type API struct {
	*api.Base
	service spec.UserService
}

var (
	healthy     = true               // Simple health flag
	version     = "0.0.1"            // App version number, set at build time with -ldflags "-X 'main.version=1.2.3'"
	buildInfo   = "No build details" // Build details, set at build time with -ldflags "-X 'main.buildInfo=Foo bar'"
	serviceName = "users"
	defaultPort = 9003
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

	// Our API wraps the Base API and UserService
	api := API{
		api.NewBase(serviceName, version, buildInfo, healthy, router),
		impl.NewService(serviceName),
	}

	// Add routes for this service
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
