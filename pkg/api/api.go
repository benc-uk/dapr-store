// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Base API that all services implement and extend
// ----------------------------------------------------------------------------

package api

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/mux"
)

// APIBase holds a standard set of values for all services & APIs
type APIBase struct {
	ServiceName string
	Healthy     bool
	Version     string
	BuildInfo   string
}

//
// AddCommonRoutes is a function
//
func (api *APIBase) AddCommonRoutes(router *mux.Router) {
	router.HandleFunc("/healthz", api.HealthCheck)
	router.HandleFunc("/api/healthz", api.HealthCheck)
	router.HandleFunc("/status", api.Status)
	router.HandleFunc("/api/status", api.Status)
}

//
// HealthCheck - Simple health check endpoint, returns 204 when healthy
//
func (api *APIBase) HealthCheck(resp http.ResponseWriter, req *http.Request) {
	if api.Healthy {
		resp.WriteHeader(http.StatusNoContent)
		return
	}
	resp.WriteHeader(http.StatusServiceUnavailable)
}

//
// Status - status information data - Remove if you like
//
func (api *APIBase) Status(resp http.ResponseWriter, req *http.Request) {
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
