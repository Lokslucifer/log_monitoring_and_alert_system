package service
import (
	"log_processor/internal/repository"
	"log_processor/internal/dto"
	"log_processor/internal/models"
)

type LogFilterService struct {
	repo repository.LogStorage

}

func NewLogFilterService(repo repository.LogStorage)*LogFilterService{
	return &LogFilterService{repo: repo}
}

func (ser *LogFilterService)FilterLogs(filter dto.FilterDTO) ([]models.LogEntry,error){
	return ser.repo.GetLogsByFilter(filter)
}

