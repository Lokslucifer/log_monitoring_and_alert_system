package main

import (
	"alert_system/internal/alert_consumer"
	"alert_system/internal/alert_sender"
	"log"
	"os"
	"github.com/joho/godotenv"
)


func main() {

	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//create a log file and setting it as default globally
	logFilePath := os.Getenv("LOG_FILE_PATH")
    if logFilePath == "" {
        logFilePath = "./log/alert_system.log" // fallback default
    }

    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer file.Close()

    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)


	slackwebhookURL := os.Getenv("SLACK_WEB_HOOK")
	if slackwebhookURL== "" {
		log.Fatal("SLACK_WEB_HOOK is not set in environment")
	}

	queue_name := os.Getenv("RABBIT_MQ_QUEUE")
	if queue_name == "" {
		log.Fatal("RABBIT_MQ_QUEUE is not set in environment")
	}

	slack_alerter:=alertsender.NewSlackAlertSender(slackwebhookURL)
	consumer,err:=alertconsumer.NewRabbitMQAlertConsumer(queue_name)
	if(err!=nil){
		log.Fatalf("error in creating alert consumer:%v",err)
	}
	consumer.StartConsumingLog(slack_alerter)
	defer consumer.Close()


}
