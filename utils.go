package main

import (
	"os"
	"os/exec"
	"strings"
)

func Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
func LowerAll(vals ...string) []string {
	res := make([]string, len(vals))
	for i, v := range vals {
		res[i] = strings.ToLower(v)
	}
	return res
}
