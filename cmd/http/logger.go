package http

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
	level string
}

func NewLogger(level string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[http-server]", log.LstdFlags|log.Lshortfile),
		level:  level,
	}
}

func (l *Logger) Info(format string, args ...any) {
	l.Printf("[INFO] "+format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.Printf("[ERROR] "+format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	l.Printf("[DEBUG] "+format, args...)
}
