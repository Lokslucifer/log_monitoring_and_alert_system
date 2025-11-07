package alertdispatcher
import (
	"log_processor/internal/models"
)

type AlertPublisher interface{
	PublishLog(logEntry models.LogEntry) error
	Close()
}