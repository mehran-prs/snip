package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	Cfg = &Config{Verbose: true}
	buf := bytes.NewBuffer(nil)
	logger.SetOutput(buf)
	Verbose("a", "b")
	assertTrue(t, strings.HasSuffix(buf.String(), ": ab\n"), "val", buf.String())

	buf.Reset()
	Cfg.Verbose = false
	Verbose("a", "b")
	assertEqual(t, buf.Len(), 0)

	buf.Reset()
	Error("c", "d")
	assertTrue(t, strings.HasSuffix(buf.String(), ": cd\n"), "val", buf.String())
}

func TestSetLoggerFile(t *testing.T) {
	Cfg = &Config{}
	// Invalid file
	err := SetLoggerFile("a/b.log")
	assertTrue(t, err != nil)

	// Opening valid file
	bytes := make([]byte, 100)
	f, err := os.CreateTemp("", "abc")
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	assertTrue(t, err == nil)
	assertTrue(t, SetLoggerFile(f.Name()) == nil)
	Error("a", "b")
	f.Read(bytes)
	assertTrue(t, strings.Contains(string(bytes), ": ab"))
	assertTrue(t, CloseLoggerFile(os.Stderr) == nil)
	// File must be closed now.
	assertTrue(t, logFile.Close() != nil)
	stderr, ok := logger.Writer().(*os.File)
	assertTrue(t, ok)
	assertEqual(t, stderr, os.Stderr)
}
