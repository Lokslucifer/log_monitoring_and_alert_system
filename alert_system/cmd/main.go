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
