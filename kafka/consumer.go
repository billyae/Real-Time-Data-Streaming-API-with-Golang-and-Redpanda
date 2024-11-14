package kafka

import (
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "data-streaming/utils"
)

func Consume(streamID string) []map[string]interface{} {
    topic := "stream_" + streamID
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": "localhost:9092",
        "group.id":          streamID,
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        utils.Logger.Fatal("Failed to create consumer:", err)
    }
    defer consumer.Close()

    consumer.SubscribeTopics([]string{topic}, nil)

    var results []map[string]interface{}
    for {
        msg, err := consumer.ReadMessage(-1)
        if err == nil {
            data := map[string]interface{}{"result": string(msg.Value)}
            results = append(results, data)
        } else {
            break
        }
    }
    return results
}
