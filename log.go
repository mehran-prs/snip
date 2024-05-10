package main

import (
	"log"
	"os"
	"sync/atomic"
)

var (
	logDebug = log.New(os.Stderr, "\x1B[36mDEBUG: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	logInfo  = log.New(os.Stderr, "\x1B[32mINFO: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	logWarn  = log.New(os.Stderr, "\x1B[35mWARN: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	logError = log.New(os.Stderr, "\x1B[31mERROR: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func LogLevelFromString(level string) int32 {
	switch level {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		return -1
	}
}

// logLevel default is warning.
var logLevel atomic.Int32

func init() {
	logLevel.Store(LogLevelWarn)
}

func SetLogLevel(level int32) { logLevel.Store(level) }

// Debug Log
func Debug(v ...interface{}) {
	if LogLevelDebug >= logLevel.Load() {
		logDebug.Println(v...)
	}
}

// Info Log
func Info(v ...interface{}) {
	if LogLevelInfo >= logLevel.Load() {
		logInfo.Println(v...)
	}
}

// Warn log
func Warn(v ...interface{}) {
	if LogLevelWarn >= logLevel.Load() {
		logWarn.Println(v...)
	}
}

// Error Log
func Error(v ...interface{}) {
	if LogLevelError >= logLevel.Load() {
		logError.Println(v...)
	}
}
