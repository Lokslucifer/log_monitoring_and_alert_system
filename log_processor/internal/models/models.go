package models

import "time"

// LogEntry represents a structured application log record
type LogEntry struct {
	Level      string    `json:"level"`       // e.g. INFO, WARN, ERROR
	Timestamp  time.Time `json:"timestamp"`   // parsed from log line
	SourceFile string    `json:"source_file"` // e.g. service.go
	LineNumber int       `json:"line_number"` // e.g. 36
	Message    string    `json:"message"`     // e.g. Email Verified successfully
}
