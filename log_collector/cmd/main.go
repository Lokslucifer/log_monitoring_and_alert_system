package main

import (
	"log_collector/internal/service"
)

func main(){
	
	file_name:="../sys.log"
	log_collector:=Service.NewLogCollector(file_name)
	log_collector.StartLogCollector()

}