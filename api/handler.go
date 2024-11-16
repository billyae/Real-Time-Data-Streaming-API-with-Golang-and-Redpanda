package api

import (
	"encoding/json"
	"net/http"
	"data-streaming/kafka"
	"data-streaming/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(100, 10) // 100 requests per second with burst size of 10

var validAPIKeys = map[string]bool{
	"billyae": true,
}

// Prometheus metrics
var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Number of requests received",
		},
		[]string{"method", "endpoint"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_duration_seconds",
			Help: "Duration of request handling in seconds",
		},
		[]string{"method", "endpoint"},
	)
)

// init registers the Prometheus metrics
func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
}

// APIKeyAuthMiddleware checks for a valid API key
func APIKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if !validAPIKeys[apiKey] {
			http.Error(w, "Forbidden: Invalid API Key", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// InitializeRoutes initializes the API routes
func InitializeRoutes(router *mux.Router) {
	apiRoutes := router.PathPrefix("/").Subrouter()
	apiRoutes.Use(APIKeyAuthMiddleware)

	apiRoutes.HandleFunc("/stream/start", StartStream).Methods("POST")
	apiRoutes.HandleFunc("/stream/{stream_id}/send", SendData).Methods("POST")
	apiRoutes.HandleFunc("/stream/{stream_id}/results", GetResults).Methods("GET")
}

// StartStream creates a new stream
func StartStream(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		return
	}

	timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, "/stream/start"))
	defer timer.ObserveDuration()
	requestsTotal.WithLabelValues(r.Method, "/stream/start").Inc()

	streamID := utils.GenerateStreamID()
	utils.Logger.Info("Started new stream with ID:", streamID)
	json.NewEncoder(w).Encode(map[string]string{"stream_id": streamID})
}

// SendData sends data to a stream
func SendData(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		return
	}

	// Measure request duration
	timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, "/stream/{stream_id}/send"))
	defer timer.ObserveDuration()
	requestsTotal.WithLabelValues(r.Method, "/stream/{stream_id}/send").Inc()

	vars := mux.Vars(r)
	streamID := vars["stream_id"]

	// Decode JSON data from request body
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.Logger.Error("Invalid data")
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	/// Send data to Kafka
	err := kafka.Produce(streamID, data)
	if err != nil {
		utils.Logger.Error("Failed to send data")
		http.Error(w, "Failed to send data", http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "data sent"})
}

// GetResults retrieves results from a stream
func GetResults(w http.ResponseWriter, r *http.Request) {

	// Check rate limit
	if !limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		return
	}

	// Measure request duration
	timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, "/stream/{stream_id}/results"))
	defer timer.ObserveDuration()
	requestsTotal.WithLabelValues(r.Method, "/stream/{stream_id}/results").Inc()

	// Get stream ID from URL
	vars := mux.Vars(r)
	streamID := vars["stream_id"]

	// Consume data from Kafka
	results := kafka.Consume(streamID)
	utils.Logger.Info("Data received from Kafka:", results)
	json.NewEncoder(w).Encode(results)
}
