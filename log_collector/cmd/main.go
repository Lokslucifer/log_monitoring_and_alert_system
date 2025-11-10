package main

import (
	"log"
	"log_collector/internal/log_streamer"
	"log_collector/internal/service"
	"os"
)

func main() {

	//create a log file and setting it as default globally
	logFilePath := os.Getenv("LOG_FILE_PATH")
    if logFilePath == "" {
        logFilePath = "./log/log_collector.log" // fallback default
    }

    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer file.Close()

    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//Geting simulated log file path and using kafka to stream the log message in that file to processor.
	simlogFilePath := os.Getenv("SIMULATED_LOG_FILE_PATH")
	if simlogFilePath == "" {
		log.Fatal("SIMULATED_LOG_FILE_PATH is not set in environment")
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

	logCollector := service.NewLogCollector(simlogFilePath, msgSender)
	logCollector.StartLogCollector()
}
