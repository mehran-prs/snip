package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func assertEqual[T comparable](t *testing.T, a, b T) {
	t.Helper()

	if a != b {
		t.Errorf(" two values are not equal. a: %v, b: %v", a, b)
	}
}

func assertEqualSlice[T ~[]E, E comparable](t *testing.T, a, b T) {
	t.Helper()

	if len(a) != len(b) {
		t.Fatalf("lengths are not equal. a: %v, b: %v", a, b)
	}

	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("elements are not equal. a: %v, b: %v", a, b)
		}
	}
}

func assertTrue(t *testing.T, isTrue bool, args ...any) {
	t.Helper()

	if !isTrue {
		t.Fatal("wanted true, got false", fmt.Sprint(args...))
	}
}

func TestCommand(t *testing.T) {
	t.Helper()
	cmd := Command("a", "b")

	assertEqualSlice(t, cmd.Args, []string{"a", "b"})
	stdin, ok := cmd.Stdin.(*os.File)
	assertTrue(t, ok)
	assertEqual(t, stdin, os.Stdin)

	stderr, ok := cmd.Stderr.(*os.File)
	assertTrue(t, ok)
	assertEqual(t, stderr, os.Stderr)

	stdout, ok := cmd.Stdout.(*os.File)
	assertTrue(t, ok)
	assertEqual(t, stdout, os.Stdout)
}

func TestDefaultStr(t *testing.T) {
	assertEqual(t, DefaultStr("", "b"), "b")
	assertEqual(t, DefaultStr("a", "b"), "a")
	assertEqual(t, DefaultStr("", "", "c"), "c")
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"normal",
			"hello world",
			[]string{"hello", "world"},
		},
		{
			"quote",
			"hello \"world hello\"",
			[]string{"hello", "world hello"},
		},
		{
			"utf-8",
			"hello 世界",
			[]string{"hello", "世界"},
		},
		{
			"space",
			"hello\\ world",
			[]string{"hello world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := parseCommand(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expect %v, got %v", tt.want, got)
			}
		})
	}
}

func TestJoinPaths(t *testing.T) {
	assertEqual(t, "a/b/", JoinPaths("a", "b/"))
	assertEqual(t, "a/b", JoinPaths("a", "b"))
}
