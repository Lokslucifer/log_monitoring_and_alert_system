package repository

import (
	"log_processor/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteStorage struct {
	db *gorm.DB
}

// NewSQLiteStorage initializes and migrates the SQLite DB
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.LogEntry{}); err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

// SaveLog inserts a log record
func (s *SQLiteStorage) AddLog(entry *models.LogEntry) error {
	return s.db.Create(entry).Error
}

// GetLogsByLevel returns all logs for a specific log level (e.g. "INFO", "ERROR")
func (s *SQLiteStorage) GetLogsByLevel(level string) ([]models.LogEntry, error) {
	var logs []models.LogEntry
	err := s.db.Where("level = ?", level).Find(&logs).Error
	return logs, err
}

// Close closes the underlying DB connection
func (s *SQLiteStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
