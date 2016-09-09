package ss

import (
	"log"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
)

var logLevel LogLevel

func init() {
	logLevel = DebugLevel
}

func SetLevel(level LogLevel) {
	logLevel = level
}

func Debugf(format string, args ...interface{}) {
	if logLevel <= DebugLevel {
		log.Printf("[DEBG] "+format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logLevel <= InfoLevel {
		log.Printf("[INFO] "+format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logLevel <= WarnLevel {
		log.Printf("[WARN] "+format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logLevel <= ErrorLevel {
		log.Printf("[ERRO] "+format, args...)
	}
}

func Panicf(format string, args ...interface{}) {
	if logLevel <= PanicLevel {
		log.Printf("[PANC] "+format, args...)
	}
}
