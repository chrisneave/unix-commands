package main

import (
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	examples := []struct {
		source string
		lines  int64
	}{
		{source: "", lines: 1},
		{source: "\n", lines: 2},
		{source: "Foo\n", lines: 2},
		{source: "Foo", lines: 1},
		{source: "Foo\nBar\n", lines: 3},
		{source: "Foo\nBar\nBaz", lines: 3},
	}

	for _, example := range examples {
		reader := strings.NewReader(example.source)
		result := count(reader)
		if result.lines != example.lines {
			t.Errorf("count(\"%s\") result.lines => %d, want %d", example.source, result.lines, example.lines)
		}
	}
}
