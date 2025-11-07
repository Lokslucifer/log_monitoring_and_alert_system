package main

import (
	"fmt"
	"log"
	"os"



	// "bufio"
	// "strings"
	alertdispatcher "log_processor/internal/alert_dispatcher"
	"sync"

	logstreamer "log_processor/internal/log_streamer"
	"log_processor/internal/repository"
	"log_processor/internal/service"
)

// func simulate(ch chan<- string) {
// 	logs := [...]string{
// 		"INFO: 2025/11/05 12:12:38 service.go:36: Email Verified successfully",
// 		"INFO: 2025/11/05 12:12:41 service.go:36: Login successful",
// 		"ERROR: 2025/11/05 12:12:44 service.go:42: Invalid Request",
// 		"ERROR: 2025/11/05 12:12:47 service.go:42: Invalid Request",
// 	}

// 	for _, log := range logs {
// 		fmt.Println(log, "- sending")
// 		ch <- log
// 	}
// 	close(ch)
// }

func main() {



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
		fmt.Printf("error initializing database: %v\n", err)
		return
	}
	queue_name := os.Getenv("RABBIT_MQ_QUEUE")
	if queue_name == "" {
		log.Fatal("RABBIT_MQ_QUEUE is not set in environment")
	}

	alertpub, err := alertdispatcher.NewRabbitMQAlertPublisher(queue_name)
	fmt.Println(alertpub, "-", err, "- initialised alert pub")
	if err != nil {
		log.Fatalf("error in creating alert published:%v", err)
	}

	logProcessor := service.NewLogProcessor(logRepo, ch, alertpub)
	go logProcessor.ProcessLog(&wg)

	// reader := bufio.NewReader(os.Stdin)

	// loop:
	// for {
	// 	// break
	// 	fmt.Println("\n==== Log Viewer ====")
	// 	fmt.Println("1. View all logs")
	// 	fmt.Println("2. View logs by level")
	// 	fmt.Println("3. Exit")
	// 	fmt.Print("Enter choice: ")

	// 	var choice int
	// 	fmt.Scanln(&choice)

	// 	switch choice {
	// 	case 1:
	// 		logs, err := logRepo.GetAllLogs() // empty means all
	// 		if err != nil {
	// 			fmt.Println("Error fetching logs:", err)
	// 			continue
	// 		}
	// 		fmt.Println("\n--- All Logs ---")
	// 		for _, logEntry := range logs {
	// 			fmt.Printf("[%s] %s\n", logEntry.Level, logEntry.Message)
	// 		}

	// 	case 2:
	// 		fmt.Print("Enter level (INFO/ERROR/WARN): ")
	// 		levelInput, _ := reader.ReadString('\n')
	// 		levelInput = strings.TrimSpace(levelInput)
	// 		logs, err := logRepo.GetLogsByLevel(levelInput)
	// 		if err != nil {
	// 			fmt.Println("Error fetching logs:", err)
	// 			continue
	// 		}
	// 		fmt.Printf("\n--- Logs with level: %s ---\n", strings.ToUpper(levelInput))
	// 		for _, logEntry := range logs {
	// 			fmt.Printf("[%s] %s\n", logEntry.Level, logEntry.Message)
	// 		}

	// 	case 3:
	// 		fmt.Println("Exiting...")
	// 		receiver.StopReceiver()
	// 		break loop

	// 	default:
	// 		fmt.Println("Invalid choice, please try again.")
	// 	}
	// }
	wg.Wait()
}
