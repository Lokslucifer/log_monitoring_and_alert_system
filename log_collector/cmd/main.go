package main

import (
	"log_collector/internal/service"
	"log_collector/internal/message_sender"
)

func main(){
	file_name:="./log/sys.log"
	topic_name:="log_processor"
	msg_sender:=message_sender.NewKafkaSender(topic_name)
	log_collector:=Service.NewLogCollector(file_name,msg_sender)
	log_collector.StartLogCollector()

}