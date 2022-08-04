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
	"path/filepath"
	"time"

	"github.com/benc-uk/dapr-store/pkg/env"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload" // Autoloads .env file if it exists
)

// Defaults, can be modified with env variables
var (
	defaultStaticPath = "./dist" // Change with env: STATIC_DIR
	defaultPort       = 8000     // Change with env: PORT
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behaviour for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("### Dapr Store: frontend host starting...")

	router := mux.NewRouter()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			log.Printf("### Problem with healthz endpoint %s\n", err)
		}
	})

	staticPath := env.GetEnvString("STATIC_DIR", defaultStaticPath)
	spa := spaHandler{
		staticPath: staticPath,
		indexPath:  "index.html",
	}

	router.HandleFunc("/config", routeConfig)
	router.PathPrefix("/").Handler(spa)

	serverPort := env.GetEnvInt("PORT", defaultPort)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", serverPort),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("### Static content path: %v\n", spa.staticPath)
	log.Printf("### Server listening on: %v\n", serverPort)
	log.Fatal(srv.ListenAndServe())
}

//
// Simple config endpoint, returns API_ENDPOINT & AUTH_CLIENT_ID vars to front end
//
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
