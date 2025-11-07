package service

import (
	"fmt"
	alertdispatcher "log_processor/internal/alert_dispatcher"
	"log_processor/internal/repository"
	"log_processor/internal/utils"
	"sync"
)

type LogProcessor struct {
	repo     repository.LogStorage
	receiver <-chan string
	alertpub alertdispatcher.AlertPublisher
}

func NewLogProcessor(repo repository.LogStorage, receiver <-chan string, alertpub alertdispatcher.AlertPublisher) *LogProcessor {

	return &LogProcessor{repo: repo, receiver: receiver,alertpub: alertpub}
}

func (lp *LogProcessor) ProcessLog(wg *sync.WaitGroup) {
	defer wg.Done()

	for logline := range lp.receiver {
		// fmt.Printf("%v",logline)
		log, err := utils.ParseLogLine(logline)

		if err != nil {
			fmt.Printf("error in parsing log:%v", err)
			continue
		}

		if log.Level == "ERROR" {
			
			err := lp.alertpub.PublishLog(*log)
			if err != nil {
				fmt.Printf("error in publishing log:%v", err)
			}

		}
		err = lp.repo.AddLog(log)
		if err != nil {
			fmt.Printf("error in adding log:%v", err)
		}
	}
}
