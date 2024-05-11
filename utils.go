package main

import (
	"os"
	"os/exec"
)

func Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func defaultVal(val string, def ...string) string {
	if val != "" {
		return val
	}
	for _, d := range def {
		if d != "" {
			return d
		}
	}
	return ""
}
