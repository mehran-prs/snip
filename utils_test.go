package main

import (
	"reflect"
	"testing"
)

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
