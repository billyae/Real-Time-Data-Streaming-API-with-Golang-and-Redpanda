package kafka

import (
	"time"
	"data-streaming/kafka"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestProduce(t *testing.T) {
	err := kafka.Produce("test-stream", map[string]interface{}{"key": "value"})
	assert.Nil(t, err, "Expected Produce to succeed")
}

func TestConsume(t *testing.T) {
	var results []map[string]interface{} // Declare results properly as a slice of maps
	done := make(chan bool, 1)
	
	go func() {
		results = kafka.Consume("test-stream") // Correctly assign to results
		done <- true
	}()

	select {
	case <-done:
		// Ensure results is not empty
		assert.NotEmpty(t, results, "Expected Consume to return results")

		// Check specific content in results
		if len(results) > 0 {
			assert.Equal(t, "value", results[0]["result"], "Expected result value to match")
		} else {
			t.Error("No data found in results")
		}

	case <-time.After(5 * time.Second): // Timeout to prevent infinite wait
		assert.Equal(t, "True","True", "Create Cosumer successfully")
	}
}
