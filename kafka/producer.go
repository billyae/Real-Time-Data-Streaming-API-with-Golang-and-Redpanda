package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"data-streaming/utils"
)

// Produce sends a message to a Kafka topic
var producer *kafka.Producer

// init creates a new Kafka producer
func init() {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		utils.Logger.Fatal("Failed to create producer:", err)
	}
}

// Produce sends a message to a Kafka topic
func Produce(streamID string, data map[string]interface{}) error {
	topic := "stream_" + streamID
	message, _ := json.Marshal(data)

	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}
