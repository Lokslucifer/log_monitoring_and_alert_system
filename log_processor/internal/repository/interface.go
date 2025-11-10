package repository

import (
	"log_processor/internal/models"
	"log_processor/internal/dto"
)

// LogStorage defines the contract for storing and retrieving logs
type LogStorage interface {
	AddLog(entry *models.LogEntry) error
	GetLogsByLevel(level string) ([]models.LogEntry, error)
	GetAllLogs()([]models.LogEntry,error)
	GetLogsByFilter(filter dto.FilterDTO)([]models.LogEntry,error)
	Close() error
}
