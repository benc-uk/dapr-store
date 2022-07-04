// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Middleware and optional extra routes available to any API
// ----------------------------------------------------------------------------

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/elastic/go-sysinfo"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Status struct {
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
	Uptime     string `json:"uptime"`
}

// AddLogging provides a route for the root URL path, if you need it
func (b *Base) AddRoot(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

// AddLogging provides logging middleware to the API in Apache Common Log Format.
func (b *Base) AddLogging(r *mux.Router) {
	r.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})
}

// AddCORS provides CORS middleware to the API
func (b *Base) AddCORS(origins []string, r *mux.Router) {
	r.Use(handlers.CORS(handlers.AllowedOrigins(origins)))
}

// AddMetrics adds Prometheus metrics to the API
func (b *Base) AddMetrics(r *mux.Router) {
	r.Handle("/metrics", promhttp.Handler())

	durationHistogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "response_duration_seconds",
		Help:        "A histogram of request latencies.",
		Buckets:     []float64{.001, .01, .1, .2, .5, 1, 2, 5},
		ConstLabels: prometheus.Labels{"handler": b.ServiceName},
	}, []string{"method"})

	r.Use(func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerDuration(durationHistogram, next)
	})
}

// AddHealth adds a health check endpoint to the API
func (b *Base) AddHealth(r *mux.Router) {
	// Add health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if b.Healthy {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Service %s is not healthy", b.ServiceName)))
		}
	})
}

// AddStatus adds a status & info endpoint to the API
func (b *Base) AddStatus(r *mux.Router) {
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		host, _ := sysinfo.Host()
		host.Info().Uptime()

		status := Status{
			Service:    b.ServiceName,
			Healthy:    b.Healthy,
			Version:    b.Version,
			BuildInfo:  b.BuildInfo,
			Hostname:   host.Info().Hostname,
			GoVersion:  runtime.Version(),
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			CPU:        runtime.NumCPU(),
			ClientAddr: r.RemoteAddr,
			ServerHost: r.Host,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(status)
	})
}
