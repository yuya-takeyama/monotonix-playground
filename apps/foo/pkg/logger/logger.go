package logger

import (
	"fmt"
	"log"
	"time"
)

type Logger struct {
	prefix string
}

func New(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(message string) {
	l.log("INFO", message)
}

func (l *Logger) Error(message string) {
	l.log("ERROR", message)
}

func (l *Logger) log(level, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	logMessage := fmt.Sprintf("[%s] %s [%s]: %s", timestamp, level, l.prefix, message)
	log.Println(logMessage)
}