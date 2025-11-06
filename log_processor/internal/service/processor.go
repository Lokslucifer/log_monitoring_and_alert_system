package service

import (
	"fmt"
	"log_processor/internal/repository"
	"log_processor/internal/utils"
	"sync"

)

type LogProcessor struct {
	repo repository.LogStorage
	receiver <-chan string

}

func NewLogProcessor(repo repository.LogStorage,receiver <-chan string)(*LogProcessor){

	return &LogProcessor{repo: repo,receiver: receiver}
}

func (lp *LogProcessor)ProcessLog(wg *sync.WaitGroup){
	defer wg.Done()
	
	for logline :=range lp.receiver{
		// fmt.Printf("%v",logline)
		log,err:=utils.ParseLogLine(logline)
		if(err!=nil){
			fmt.Printf("error in parsing log:%v",err)
			continue
		}
		err=lp.repo.AddLog(log)
		if(err!=nil){
			fmt.Printf("error in adding log:%v",err)
		}
	}
}


