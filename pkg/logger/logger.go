package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)


type LogLevel int

const (

	DEBUG LogLevel = iota
	
	INFO
	
	WARN

	ERROR

	FATAL
)


type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, err interface{})
	Fatal(msg string, err interface{})
}


type SimpleLogger struct {
	level  LogLevel
	logger *log.Logger
}


func NewLogger() Logger {
	level := INFO
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		level = DEBUG
	}

	return &SimpleLogger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}


func (l *SimpleLogger) formatLog(level, msg string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	
	var argsStr string
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			argsStr += fmt.Sprintf(" %v=%v", args[i], args[i+1])
		}
	}
	
	return fmt.Sprintf("[%s] [%s] %s%s", timestamp, level, msg, argsStr)
}


func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Println(l.formatLog("DEBUG", msg, args...))
	}
}


func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if l.level <= INFO {
		l.logger.Println(l.formatLog("INFO", msg, args...))
	}
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if l.level <= WARN {
		l.logger.Println(l.formatLog("WARN", msg, args...))
	}
}


func (l *SimpleLogger) Error(msg string, err interface{}) {
	if l.level <= ERROR {
		l.logger.Println(l.formatLog("ERROR", msg, "error", err))
	}
}


func (l *SimpleLogger) Fatal(msg string, err interface{}) {
	l.logger.Println(l.formatLog("FATAL", msg, "error", err))
	os.Exit(1)
}