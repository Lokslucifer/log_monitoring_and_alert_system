package main

import (
	"log_processor/internal/message_receiver"
	"log_processor/internal/service"
	"fmt"
)

func main(){
	
	topic_name:="log_processor"
	ch:=make(chan string)
	receiver:=message_receiver.NewKafkaReceiver(topic_name)
	receiver.ReceiveMessage(ch)
	return
	dbPath := "/data/logs.db"
	log_processor:=Service.NewLogProcessor(fmt.Sprintf("sqlite://%s", dbPath),ch)
	log_processor.ProcessLog()

}