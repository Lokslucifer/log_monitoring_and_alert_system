package main

import (
	"fmt"
	"log"
	"os"

	alertdispatcher "log_processor/internal/alert_dispatcher"
	"sync"

	logstreamer "log_processor/internal/log_streamer"
	"log_processor/internal/repository"
	"log_processor/internal/service"
	"log_processor/internal/handler/v1"
	"github.com/gin-gonic/gin"
)



func main() {


	logFilePath := os.Getenv("LOG_FILE_PATH")
    if logFilePath == "" {
        logFilePath = "./log/log_processor.log" // fallback default
    }

    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer file.Close()

    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)



	topicName := os.Getenv("KAFKA_TOPIC")
	if topicName == "" {
		log.Fatal("KAFKA_TOPIC is not set in environment")
	}
	log_consumer, err := logstreamer.NewKafkaLogConsumer(topicName)
	if err != nil {
		log.Fatalf("error in creating kafka log consumer:%v", err)

	}
	ch := make(chan string)
	var wg sync.WaitGroup

	wg.Add(2)
	go log_consumer.StartConsuming(ch, &wg)
	db_directory := os.Getenv("DB_DIRECTORY")
	if db_directory == "" {
		log.Fatal("DB_DIRECTORY is not set in environment")
	}

	os.MkdirAll(db_directory, os.ModePerm)
	dbPath := fmt.Sprintf("%v/logs.db", db_directory)

	logRepo, err := repository.NewSQLiteStorage(dbPath)
	if err != nil {
		log.Fatalf("error initializing database: %v\n", err)
		return
	}
	queue_name := os.Getenv("RABBIT_MQ_QUEUE")
	if queue_name == "" {
		log.Fatal("RABBIT_MQ_QUEUE is not set in environment")
	}

	alertpub, err := alertdispatcher.NewRabbitMQAlertPublisher(queue_name)
	if err != nil {
		log.Fatalf("error in creating alert published:%v", err)
	}

	logProcessor := service.NewLogProcessor(logRepo, ch, alertpub)
	go logProcessor.ProcessLog(&wg)

	log_filter_service:=service.NewLogFilterService(logRepo)

	handler:=v1.NewHandler(log_filter_service)



	// implement gin server for querying log database and stoping the server based on user request
	r := gin.Default()
	r.GET("/logs",handler.FilterLogsHandler)
	r.Run(":8090")

	wg.Wait()
}
