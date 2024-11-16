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

func InitializeRoutes(router *mux.Router) {
	apiRoutes := router.PathPrefix("/").Subrouter()
	apiRoutes.Use(APIKeyAuthMiddleware)

	apiRoutes.HandleFunc("/stream/start", StartStream).Methods("POST")
	apiRoutes.HandleFunc("/stream/{stream_id}/send", SendData).Methods("POST")
	apiRoutes.HandleFunc("/stream/{stream_id}/results", GetResults).Methods("GET")
}

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

func SendData(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		return
	}

	timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, "/stream/{stream_id}/send"))
	defer timer.ObserveDuration()
	requestsTotal.WithLabelValues(r.Method, "/stream/{stream_id}/send").Inc()

	vars := mux.Vars(r)
	streamID := vars["stream_id"]

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.Logger.Error("Invalid data")
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	err := kafka.Produce(streamID, data)
	if err != nil {
		utils.Logger.Error("Failed to send data")
		http.Error(w, "Failed to send data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "data sent"})
}

func GetResults(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
		return
	}

	timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, "/stream/{stream_id}/results"))
	defer timer.ObserveDuration()
	requestsTotal.WithLabelValues(r.Method, "/stream/{stream_id}/results").Inc()

	vars := mux.Vars(r)
	streamID := vars["stream_id"]

	results := kafka.Consume(streamID)
	utils.Logger.Info("Data received from Kafka:", results)
	json.NewEncoder(w).Encode(results)
}
