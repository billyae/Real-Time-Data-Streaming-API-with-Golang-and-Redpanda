package api

import (
    "encoding/json"
    "net/http"
    // "github.com/confluentinc/confluent-kafka-go/kafka"
    localKafka "data-streaming/kafka"
    "data-streaming/utils"
    "github.com/gorilla/mux"
)

func InitializeRoutes(router *mux.Router) {
    router.HandleFunc("/stream/start", StartStream).Methods("POST")
    router.HandleFunc("/stream/{stream_id}/send", SendData).Methods("POST")
    router.HandleFunc("/stream/{stream_id}/results", GetResults).Methods("GET")
}

func StartStream(w http.ResponseWriter, r *http.Request) {
    // Generate a unique stream ID
    streamID := utils.GenerateStreamID()
    utils.Logger.Info("Started new stream with ID:", streamID)
    json.NewEncoder(w).Encode(map[string]string{"stream_id": streamID})
}

func SendData(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    streamID := vars["stream_id"]

    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid data", http.StatusBadRequest)
        return
    }

    // Send data to Kafka
    err := localKafka.Produce(streamID, data)
    if err != nil {
        http.Error(w, "Failed to send data", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "data sent"})
}

func GetResults(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    streamID := vars["stream_id"]

    // Fetch results from Kafka consumer
    results := localKafka.Consume(streamID)
    json.NewEncoder(w).Encode(results)
}
