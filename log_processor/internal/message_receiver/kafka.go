package message_receiver

import (
	"context"
	"fmt"
	"log"
	"sync"
    "os"
	"github.com/segmentio/kafka-go"
)

type KafkaReceiver struct {
	topicName string
	reader    *kafka.Reader
	cancel    context.CancelFunc
}

func (kr *KafkaReceiver) ReceiveMessage(ch chan<- string, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("closing channel from receiver")
		close(ch)
		wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	kr.cancel = cancel

	for {
		msg, err := kr.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("Kafka reader stopped (context canceled)")
				return
			}
			fmt.Printf("Consumer error: %v\n", err)
			continue
		}
		fmt.Printf("Received log: %s\n", msg.Value)
		ch <- string(msg.Value)
	}
}

func (kr *KafkaReceiver) StopReceiver() {
	fmt.Println("Stopping receiver...")
	if kr.cancel != nil {
		kr.cancel()
	}
	kr.reader.Close()
}

func NewKafkaReceiver(topicName string) *KafkaReceiver {
    brokers := []string{os.Getenv("KAFKA_BROKER")}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topicName,
		GroupID:  "log_processor_group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Printf("Kafka consumer initialized for topic: %s\n", topicName)

	return &KafkaReceiver{
		topicName: topicName,
		reader:    reader,
	}
}
