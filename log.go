package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

func Verbose(v ...interface{}) {
	if Cfg.Verbose {
		_ = logger.Output(3, fmt.Sprint(v...))
	}
}

// Error Log
func Error(v ...interface{}) {
	_ = logger.Output(3, "ERROR: "+fmt.Sprint(v...))
}

// --------------------------------
// File Output
// --------------------------------
var (
	logFile   *os.File
	logFileMu sync.Mutex
)

func setLoggerFile(filename string) error {
	logFileMu.Lock()
	defer logFileMu.Unlock()
	oldFile := logFile

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	logFile = f
	logger.SetOutput(io.MultiWriter(f, os.Stderr))

	if oldFile != nil {
		return oldFile.Close()
	}

	return nil
}

func CloseLoggerFile(w io.Writer) error {
	logFileMu.Lock()
	defer logFileMu.Unlock()

	if logFile == nil {
		return nil
	}

	logger.SetOutput(w)
	return logFile.Close()
}
