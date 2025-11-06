package main

import (
	"fmt"
	"os"

	// "bufio"
	// "strings"
	"sync"

	"log_processor/internal/message_receiver"
	"log_processor/internal/repository"
	"log_processor/internal/service"
)

func simulate(ch chan<- string) {
	logs := [...]string{
		"INFO: 2025/11/05 12:12:38 service.go:36: Email Verified successfully",
		"INFO: 2025/11/05 12:12:41 service.go:36: Login successful",
		"ERROR: 2025/11/05 12:12:44 service.go:42: Invalid Request",
		"ERROR: 2025/11/05 12:12:47 service.go:42: Invalid Request",
	}

	for _, log := range logs {
		fmt.Println(log, "- sending")
		ch <- log
	}
	close(ch)
}

func main() {
	topic_name := "log_processor"
	ch := make(chan string)
	var wg sync.WaitGroup
	receiver := message_receiver.NewKafkaReceiver(topic_name)
	wg.Add(2)
	go receiver.ReceiveMessage(ch, &wg)

	os.MkdirAll("./data", os.ModePerm)
	dbPath := "./data/logs.db"

	logRepo, err := repository.NewSQLiteStorage(dbPath)
	if err != nil {
		fmt.Printf("error initializing database: %v\n", err)
		return
	}

	logProcessor := service.NewLogProcessor(logRepo, ch)
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
