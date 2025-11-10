package repository

import (
	"log_processor/internal/dto"
	"log_processor/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteStorage struct {
	db *gorm.DB
}

// NewSQLiteStorage initializes and migrates the SQLite DB
func NewSQLiteStorage(dbPath string) (LogStorage, error) {
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

func (s *SQLiteStorage) GetAllLogs() ([]models.LogEntry, error) { 
	var logs []models.LogEntry 
	err := s.db.Find(&logs).Error 
	return logs, err 
}

func (s *SQLiteStorage) GetLogsByFilter(filter dto.FilterDTO) ([]models.LogEntry, error) {
	var logs []models.LogEntry

	query := s.db.Model(&models.LogEntry{})

	// Filter by log levels if provided
	if len(filter.Levels) > 0 {
		query = query.Where("level IN ?", filter.Levels)
	}

	// Filter by search text in message or source (example fields)
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("LOWER(message) LIKE ? OR LOWER(source_file) LIKE ?", searchPattern, searchPattern)
	}

	// Filter by time range
	if !filter.From.IsZero() {
		query = query.Where("timestamp >= ?", filter.From)
	}

	if !filter.To.IsZero() {
		query = query.Where("timestamp <= ?", filter.To)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Order by most recent logs first
	query = query.Order("timestamp DESC")

	// Execute query
	err := query.Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Close closes the underlying DB connection
func (s *SQLiteStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}


