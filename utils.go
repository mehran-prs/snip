package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func DefaultStr(val string, def ...string) string {
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

func parseCommand(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for _, c := range command {

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

// JoinPaths is just like path.Join method, but doesn't remove the last path separator from the joined paths.
// e.g., JoinPaths("a","b/") returns "a/b/" insted of "a/b"
func JoinPaths(elem ...string) string {
	lastElem := elem[len(elem)-1]
	appenPathSeparator := len(lastElem) != 0 && lastElem[len(lastElem)-1] == os.PathSeparator

	res := path.Join(elem...)
	if appenPathSeparator {
		res = res + string(os.PathSeparator)
	}
	return res
}
