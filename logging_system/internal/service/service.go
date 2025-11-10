package Service

import (
	"log"
	"os"
)

// Logger struct to encapsulate log behavior
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

// NewLogger creates a new logger instance
func NewLogger(logFile string) (*Logger, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}

	info := log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errLog := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		infoLogger:  info,
		errorLogger: errLog,
		file:        file,
	}, nil
}

// Info logs general information
func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)

}

// Error logs errors
func (l *Logger) Error(err error) {
	l.errorLogger.Println(err)

}

// Close closes the log file
func (l *Logger) Close() {
	l.file.Close()
}
