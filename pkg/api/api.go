// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Base API that all services implement and extend
// ----------------------------------------------------------------------------

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

// Base holds a standard set of values for all services & APIs
type Base struct {
	ServiceName string
	Healthy     bool
	Version     string
	BuildInfo   string
}

//
// New creates and returns a new Base API instance
//
func NewBase(name, ver, info string, healthy bool, router *mux.Router) *Base {
	b := &Base{
		ServiceName: name,
		Healthy:     healthy,
		Version:     ver,
		BuildInfo:   info,
	}
	router.HandleFunc("/healthz", b.HealthCheck)
	router.HandleFunc("/api/healthz", b.HealthCheck)
	router.HandleFunc("/status", b.Status)
	router.HandleFunc("/api/status", b.Status)

	// Add middleware for logging
	router.Use(b.loggingMiddleware)
	return b
}

//
// HealthCheck - Simple health check endpoint, returns 204 when healthy
//
func (api *Base) HealthCheck(resp http.ResponseWriter, req *http.Request) {
	if api.Healthy {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.WriteHeader(http.StatusServiceUnavailable)
}

//
// Status - status information data - Remove if you like
//
func (api *Base) Status(resp http.ResponseWriter, req *http.Request) {
	type status struct {
		Service    string `json:"service"`
		Healthy    bool   `json:"healthy"`
		Version    string `json:"version"`
		BuildInfo  string `json:"buildInfo"`
		Hostname   string `json:"hostname"`
		OS         string `json:"os"`
		Arch       string `json:"architecture"`
		CPU        int    `json:"cpuCount"`
		GoVersion  string `json:"goVersion"`
		ClientAddr string `json:"clientAddress"`
		ServerHost string `json:"serverHost"`
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "hostname not available"
	}

	currentStatus := status{
		Service:    api.ServiceName,
		Healthy:    api.Healthy,
		Version:    api.Version,
		BuildInfo:  api.BuildInfo,
		Hostname:   hostname,
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		CPU:        runtime.NumCPU(),
		ClientAddr: req.RemoteAddr,
		ServerHost: req.Host,
	}

	statusJSON, err := json.Marshal(currentStatus)
	if err != nil {
		http.Error(resp, "Failed to get status", http.StatusInternalServerError)
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.Write(statusJSON)
}

//
// Standard CORS settings, no longer used, left for prosperity
//
/*func (api *APIBase) corsMiddleware(handler http.Handler) http.Handler {
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	return handlers.CORS(corsOrigins, corsMethods)(handler)
}*/

//
// Basic request logging,
//
func (api *Base) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Invented header to not log this request, lets us ignore things like k8s probes
		noLog := r.Header.Get("No-Log")
		if noLog != "" {
			next.ServeHTTP(w, r)
			return
		}

		// Really simple request logging
		log.Printf("### %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
