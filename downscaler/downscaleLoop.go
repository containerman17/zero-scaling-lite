package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	newNamespacesCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "new_namespaces",
		Help: "Namespaces added (including startup)",
	})
	secretCreatedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "secret_created",
		Help: "New secrets created",
	})
	errorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "errors",
		Help: "Errors (non-fatal)",
	})
)

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Panic("Error running metrics server")
	}
}
