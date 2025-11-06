package repository

import "log_processor/internal/models"

// LogStorage defines the contract for storing and retrieving logs
type LogStorage interface {
	AddLog(entry *models.LogEntry) error
	GetLogsByLevel(level string) ([]models.LogEntry, error)
	GetAllLogs()([]models.LogEntry,error)
	Close() error
}
