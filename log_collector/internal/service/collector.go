package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log_collector/internal/customErrors"
	"os"
	"time"
	"log_collector/internal/log_streamer"
	"log"
)

type LogCollector struct {
	file_name string
	rotated   chan bool
	lines     chan string
	stop      chan bool
	sender     log_streamer.LogProducer
}

func NewLogCollector(file_name string,sender log_streamer.LogProducer) *LogCollector {

	return &LogCollector{file_name: file_name, rotated: make(chan bool), lines: make(chan string), stop: make(chan bool),sender:sender}

}


func (lc *LogCollector) StartLogCollector() {
	go func(){
		for {
			err:=lc.Tailer()
			if(err!=nil){
				if(errors.Is(err,Customerrors.FileRotatedError)){
					continue
				}else{
					log.Println("Tailer error:", err)
                	lc.StopLogCollector()
					return
				}
			}
		}
	}()

	go lc.Watcher()

	for line:=range lc.lines{
		err:=lc.sender.SendLog(line)
		if(err!=nil){
			log.Println(err)
			time.Sleep(500*time.Millisecond)
		}	
	}

}

func (lc *LogCollector) StopLogCollector() {
select {
	case <-lc.stop:
		// already closed
		return
	default:
		lc.stop<-true
		close(lc.stop)
	}

}




func (lc *LogCollector) Watcher() error {
	lastStat,err := os.Stat(lc.file_name)
	if(err!=nil){
		
		// lc.stop<-true
		// close(lc.stop)
		lc.StopLogCollector()

		return fmt.Errorf("error in getting file info:%v",err)
	}
	for {
		select{
		case<-lc.stop:
			close(lc.rotated)
			log.Println("Stoping File Watcher")
			return nil
		
		default:
			currentStat,err:=os.Stat((lc.file_name))
			if(err!=nil){

				lc.StopLogCollector()
				return fmt.Errorf("error in getting file info:%v",err)
			}
			if(!os.SameFile(lastStat,currentStat)){
				lc.rotated<-true
			}
		}

		time.Sleep(1*time.Second)
	
	}

}


func (lc *LogCollector) Tailer()error {
	file,err:=os.Open(lc.file_name)
	if(err!=nil){

		lc.StopLogCollector()

		return fmt.Errorf("error in Opening File:%v",err)
	}
	defer file.Close()

	file.Seek(0,io.SeekEnd)
	reader:=bufio.NewReader(file)
	var line string

	for {
		select {
		case<-lc.stop:
			close(lc.lines)
			log.Println("Stopping Tailer")
			return nil
		
		case<-lc.rotated:
			return Customerrors.FileRotatedError
		
		default:
			line,err=reader.ReadString('\n')
			if(err!=nil){
				if(err==io.EOF){
				time.Sleep(500*time.Millisecond)
				continue
				}
				lc.StopLogCollector()
				return fmt.Errorf("error in reading log file:%v",err)
			}
			lc.lines<-line

		}


	}
}



