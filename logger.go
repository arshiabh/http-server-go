package main

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

var logger *Logger

func init() {
	logger = &Logger{
		Logger: log.New(os.Stdout, "[http-server]", log.LstdFlags|log.Lshortfile),
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
