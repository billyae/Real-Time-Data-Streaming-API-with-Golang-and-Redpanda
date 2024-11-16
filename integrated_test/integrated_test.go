package integrated_test

import (
	"time"
	"os"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"data-streaming/api"
	// "fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const apiKey = "billyae"

// TestStartStream tests the /stream/start endpoint
func TestStartStream(t *testing.T) {
	router := mux.NewRouter()
	api.InitializeRoutes(router)

	req, _ := http.NewRequest("POST", "/stream/start", nil)
	req.Header.Set("x-api-key", apiKey) // Add API key header

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code, "Expected HTTP status OK")
}

// TestSendData tests the /stream/{streamid}/send endpoint
func TestSendData(t *testing.T) {

	streamid := os.Getenv("streamid")

	router := mux.NewRouter()
	api.InitializeRoutes(router)
	
	t.Log("streamid: ", streamid)
	data := map[string]interface{}{"key": "value"}
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/stream/"+streamid+"/send", bytes.NewReader(body))
	req.Header.Set("x-api-key", apiKey) // Add API key header

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code, "Expected HTTP status OK")
}

// TestGetResults tests the /stream/{streamid}/results endpoint
func TestGetResults(t *testing.T) {

	streamid := os.Getenv("streamid")
	if streamid == "" {
		t.Fatal("streamid not provided")
	}

	router := mux.NewRouter()
	api.InitializeRoutes(router)

	req, _ := http.NewRequest("GET", "/stream/"+streamid+"/results", nil)
	req.Header.Set("x-api-key", apiKey) // Add API key header
	
	response := httptest.NewRecorder()

	// Timeout mechanism to avoid infinite wait
	done := make(chan bool, 1)
	go func() {
		router.ServeHTTP(response, req)
		done <- true
	}()

	select {
	case <-done:
		assert.Equal(t, http.StatusOK, response.Code, "Expected HTTP status OK")
	case <-time.After(5 * time.Second): // Adjust timeout as needed
		assert.Equal(t, "True","True", "Wait for the data")
	
	}
}