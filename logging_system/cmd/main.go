package main

import (
	"fmt"
	"logging_system/internal/service"
	"math/rand"
	"time"
	"os"
	"log"

)

func simulate_logging(){

	logFilePath := os.Getenv("SIMULATED_LOG_FILE_PATH")
	if logFilePath == "" {
		log.Fatal("LOG_FILE_PATH is not set in environment")
	}

	logger,err:=Service.NewLogger(logFilePath)
	if(err!=nil){
		log.Fatal("Error During Logger creation:",err)
	}
	error_messages:=[...]string{"Internal Server Error","Invalid Input","Invalid Request","Not Authorised"};
	info_messages:=[...]string{"Email Verified successfully","Login successful","Signup unsuccessful"}
	log_types:=[...]string{"ERROR","INFO"}
	var log_type string
	var log_message string

	for {
		log_type=log_types[rand.Intn(len(log_types))]
		if (log_type=="ERROR"){

			log_message=error_messages[rand.Intn(len(error_messages))]
			logger.Error(fmt.Errorf(log_message))

		}else{

			log_message=info_messages[rand.Intn(len(info_messages))]
			logger.Info(log_message)
		}
		time.Sleep(3*time.Second)

	}

}
func main() {
    logFilePath := os.Getenv("LOG_FILE_PATH")
    if logFilePath == "" {
        logFilePath = "./log/system_app.log" // fallback default
    }

    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error opening log file: %v", err)
    }
    defer file.Close()

    log.SetOutput(file)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

    simulate_logging()
}
