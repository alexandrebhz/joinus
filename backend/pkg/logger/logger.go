package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

type SimpleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

func NewLogger() Logger {
	return &SimpleLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLogger.Printf(msg, args...)
	} else {
		l.infoLogger.Println(msg)
	}
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.errorLogger.Printf(msg, args...)
	} else {
		l.errorLogger.Println(msg)
	}
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.debugLogger.Printf(msg, args...)
	} else {
		l.debugLogger.Println(msg)
	}
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.warnLogger.Printf(msg, args...)
	} else {
		l.warnLogger.Println(msg)
	}
}



