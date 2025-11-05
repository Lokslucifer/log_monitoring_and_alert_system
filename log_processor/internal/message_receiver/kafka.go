package message_receiver

import (
    "context"
    "fmt"
    "log"
    "github.com/segmentio/kafka-go"
)

// KafkaReceiver handles reading messages from Kafka
type KafkaReceiver struct {
    topicName string
    reader    *kafka.Reader
}

// ReceiveMessage continuously reads and prints messages
func (kr *KafkaReceiver) ReceiveMessage(ch chan<-string ) error {
    for {
        msg, err := kr.reader.ReadMessage(context.Background())
        if err != nil {
            fmt.Printf("Consumer error: %v\n", err)
            continue
        }
        ch<-string(msg.Value)
        // fmt.Printf("Received log: %s\n", string(msg.Value))
    }
}

// NewKafkaReceiver initializes a KafkaReceiver with a Kafka reader
func NewKafkaReceiver(topicName string) *KafkaReceiver {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  []string{"localhost:9092"},
        Topic:    topicName,
        GroupID:  "log_processor_group",
        MinBytes: 10e3,  // 10KB
        MaxBytes: 10e6,  // 10MB
    })

    log.Printf("Kafka consumer initialized for topic: %s\n", topicName)

    return &KafkaReceiver{
        topicName: topicName,
        reader:    reader,
    }
}
