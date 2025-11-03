package main

import (
	"fmt"
	"logging_system/internal/service"
	"math/rand"
	"time"
)

func simulate_logging(){

	log_file_name:="../sys.log"
	logger,err:=Service.NewLogger(log_file_name)
	if(err!=nil){
		fmt.Println("Error During Logger creation:",err)
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
func main(){
	simulate_logging()
}