// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Static HTML/ content host for frontend
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/benc-uk/go-rest-api/pkg/env"
	"github.com/benc-uk/go-rest-api/pkg/static"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// Defaults, can be modified with env variables
var (
	defaultStaticPath = "./dist" // Change with env: STATIC_DIR
	defaultPort       = 8000     // Change with env: PORT
)

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("### Dapr Store: frontend host starting...")

	router := chi.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			log.Printf("### Problem with health endpoint %s\n", err)
		}
	})

	staticPath := env.GetEnvString("STATIC_DIR", defaultStaticPath)
	spa := static.SpaHandler{
		StaticPath: staticPath,
		IndexFile:  "index.html",
	}

	router.HandleFunc("/config", routeConfig)
	router.Handle("/*", spa)

	// This stops Dapr from trying to subscribe to this app, which generated warnings
	router.HandleFunc("/dapr/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	serverPort := env.GetEnvInt("PORT", defaultPort)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("### Static content path: %v\n", spa.StaticPath)
	log.Printf("### Server listening on: %v\n", serverPort)
	log.Fatal(srv.ListenAndServe())
}

// Simple config endpoint, returns API_ENDPOINT & AUTH_CLIENT_ID vars to front end
func routeConfig(resp http.ResponseWriter, req *http.Request) {
	config := map[string]string{
		"API_ENDPOINT":   env.GetEnvString("API_ENDPOINT", "/"),
		"AUTH_CLIENT_ID": env.GetEnvString("AUTH_CLIENT_ID", ""),
	}

	configJSON, _ := json.Marshal(config)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Content-Type", "application/json")
	_, _ = resp.Write([]byte(configJSON))
}
