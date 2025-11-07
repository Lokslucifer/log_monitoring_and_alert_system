package main

import (
	"log"
	"log_collector/internal/log_streamer"
	"log_collector/internal/service"
	"os"
)

func main() {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		log.Fatal("LOG_FILE_PATH is not set in environment")
	}

	topicName := os.Getenv("KAFKA_TOPIC")
	if topicName == "" {
		log.Fatal("KAFKA_TOPIC is not set in environment")
	}

	msgSender, err := log_streamer.NewKafkaLogProducer(topicName)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer msgSender.Close() // Always good to close the producer

	logCollector := service.NewLogCollector(logFilePath, msgSender)
	logCollector.StartLogCollector()
}
