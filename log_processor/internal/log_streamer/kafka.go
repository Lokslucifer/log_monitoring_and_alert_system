package logstreamer

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaLogConsumer struct {
	topicName   string
	kafkaReader *kafka.Reader
	ctxCancel   context.CancelFunc
}

func (k *KafkaLogConsumer) StartConsuming(ch chan<- string, wg *sync.WaitGroup) {
	defer func() {
		log.Println("closing log channel")
		close(ch)
		wg.Done()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	k.ctxCancel = cancel

	for {
		msg, err := k.kafkaReader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Kafka reader stopped (context canceled)")
				return
			}
			log.Printf("Kafka consume error: %v\n", err)
			continue
		}
		log.Printf("Received log: %s\n", msg.Value)
		ch <- string(msg.Value)
	}
}

func (k *KafkaLogConsumer) Stop() {
	log.Println("Stopping KafkaLogConsumer...")
	if k.ctxCancel != nil {
		k.ctxCancel()
	}
	k.kafkaReader.Close()
}

func NewKafkaLogConsumer(topicName string) (LogConsumer, error) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092" // default fallback
	}
	log.Println(broker)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topicName,
		GroupID:  "log_processor_group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Printf("Kafka consumer initialized for topic: %s\n", topicName)

	return &KafkaLogConsumer{
		topicName:   topicName,
		kafkaReader: reader,
	}, nil
}
