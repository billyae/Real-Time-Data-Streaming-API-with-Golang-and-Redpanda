package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "data-streaming/api"
    "fmt"
)

func main() {
    // Create or open a log file
    
    router := mux.NewRouter()
    api.InitializeRoutes(router)
    
    // Expose Prometheus metrics
    router.Handle("/metrics", promhttp.Handler())

    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", router)

    
 
}
