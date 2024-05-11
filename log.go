package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// logger is the global logger.
var (
	logger   = MustNewLogger("", LogLevelWarn)
	loggerMu sync.RWMutex
)

// SetGlobalLogger sets the global logger.
func SetGlobalLogger(l *Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()
	logger = l
}

func GlobalLogger() *Logger {
	loggerMu.RLock()
	defer loggerMu.RUnlock()
	return logger
}

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

type Logger struct {
	f   *os.File // f is immutable.
	l   *log.Logger
	lvl int32
}

func NewLogger(filename string, lvl int32) (*Logger, error) {
	logger := &Logger{
		l:   log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile),
		lvl: lvl,
	}

	if filename != "" {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}

		logger.f = f
		logger.l.SetOutput(io.MultiWriter(f, os.Stderr))
	}
	return logger, nil
}

func MustNewLogger(filename string, lvl int32) *Logger {
	l, err := NewLogger(filename, lvl)
	if err != nil {
		panic(err)
	}
	return l
}

// Debug Log
func (l *Logger) Debug(v ...interface{}) {
	if LogLevelDebug >= l.lvl {
		_ = l.l.Output(3, "DEBUG: "+fmt.Sprint(v...))
	}
}

// Info Log
func (l *Logger) Info(v ...interface{}) {
	if LogLevelInfo >= l.lvl {
		_ = l.l.Output(3, "INFO: "+fmt.Sprint(v...))
	}
}

// Warn log
func (l *Logger) Warn(v ...interface{}) {
	if LogLevelWarn >= l.lvl {
		_ = l.l.Output(3, "WARN: "+fmt.Sprint(v...))
	}
}

// Error Log
func (l *Logger) Error(v ...interface{}) {
	if LogLevelError >= l.lvl {
		_ = l.l.Output(3, "ERROR: "+fmt.Sprint(v...))
	}
}

func (l *Logger) Shutdown() error {
	l.l.SetOutput(os.Stderr)
	if l.f != nil {
		return l.f.Close()
	}
	return nil
}

func Debug(v ...any) {
	GlobalLogger().Debug(v...)
}
func Info(v ...any) {
	GlobalLogger().Info(v...)
}
func Warn(v ...any) {
	GlobalLogger().Warn(v...)
}
func Error(v ...any) {
	GlobalLogger().Error(v...)
}
