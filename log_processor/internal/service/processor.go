package Service

import (
	"fmt"
	"log_processor/internal/repository"
	"log_processor/internal/utils"
)

type LogProcessor struct {
	repo repository.LogStorage
	receiver <-chan string

}

func NewLogProcessor(dburl string,receiver <-chan string)(*LogProcessor){
	repo,err:=repository.NewSQLiteStorage(dburl)
	if(err!=nil){
		fmt.Printf("error in database initiation:%v",err)
	}
	return &LogProcessor{repo: repo,receiver: receiver}
}

func (lp *LogProcessor)ProcessLog(){
	for logline :=range lp.receiver{
		log,err:=utils.ParseLogLine(logline)
		if(err!=nil){
			fmt.Printf("error in parsing log:%v",err)
		}
		err=lp.repo.AddLog(log)
		if(err!=nil){
			fmt.Printf("error in adding log:%v",err)
		}
	}
}


