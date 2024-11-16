package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"data-streaming/api"
	"fmt"
)

// main starts the server
func main() {

	// Initialize the API routes
	router := mux.NewRouter()
	api.InitializeRoutes(router)

	// Expose Prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	// Start the server
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
