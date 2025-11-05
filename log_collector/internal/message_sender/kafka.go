package message_sender

import (
    "context"
    "fmt"
    "log"

    "github.com/segmentio/kafka-go"
)

type KafkaSender struct {
    topicName string
    writer    *kafka.Writer
}

// SendMessage sends a message to the Kafka topic
func (ks *KafkaSender) SendMessage(msg string) error {
    err := ks.writer.WriteMessages(context.Background(),
        kafka.Message{
            Value: []byte(msg),
        },
    )

    if err != nil {
        return fmt.Errorf("failed to deliver message: %v", err)
    }

    fmt.Println("Message delivered successfully")
    return nil
}

// NewKafkaSender initializes and returns a KafkaSender instance
func NewKafkaSender(topicName string) *KafkaSender {
    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  []string{"localhost:9092"},
        Topic:    topicName,
        Balancer: &kafka.LeastBytes{}, // balances load across partitions
    })

    log.Printf("Kafka producer initialized for topic: %s\n", topicName)

    return &KafkaSender{
        topicName: topicName,
        writer:    writer,
    }
}
