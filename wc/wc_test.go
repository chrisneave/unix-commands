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

func TestCountWords(t *testing.T) {
	examples := []struct {
		source string
		words  int64
	}{
		{source: "", words: 0},
		{source: "Foo", words: 1},
		{source: "just three words", words: 3},
		{source: "Foo\nBar", words: 2},
		{source: "Foo\nBar\n", words: 2},
		{source: "Foo\nBar\nBaz", words: 3},
	}

	for _, example := range examples {
		reader := strings.NewReader(example.source)
		result := count(reader)
		if result.words != example.words {
			t.Errorf("count(\"%s\") result.words => %d, want %d", example.source, result.words, example.words)
		}
	}
}

func TestCountBytes(t *testing.T) {
	examples := []struct {
		source string
		bytes  int64
	}{
		{source: "", bytes: 0},
		{source: "Foo", bytes: 4},
		{source: "just three bytes", bytes: 17},
		{source: "Foo\nBar", bytes: 8},
		{source: "Foo\nBar\n", bytes: 9},
		{source: "Foo\nBar\nBaz", bytes: 12},
	}

	for _, example := range examples {
		reader := strings.NewReader(example.source)
		result := count(reader)
		if result.bytes != example.bytes {
			t.Errorf("count(\"%s\") result.bytes => %d, want %d", example.source, result.bytes, example.bytes)
		}
	}
}
