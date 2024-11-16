package kafka

import (
    
    "encoding/json"
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "data-streaming/utils"
)

var producer *kafka.Producer

func init() {
    var err error
    producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
    if err != nil {
        utils.Logger.Fatal("Failed to create producer:", err)
    }
}

func Produce(streamID string, data map[string]interface{}) error {
    topic := "stream_" + streamID
    message, _ := json.Marshal(data)

    return producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          message,
    }, nil)
}
