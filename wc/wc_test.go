package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestCountLines(t *testing.T) {
	examples := []struct {
		source string
		lines  int64
	}{
		{source: "", lines: 0},
		{source: "\n", lines: 1},
		{source: "Foo\n", lines: 1},
		{source: "Foo", lines: 1},
		{source: "Foo\nBar\n", lines: 2},
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
		{source: "Foo", bytes: 3},
		{source: "just three bytes", bytes: 16},
		{source: "Foo\nBar", bytes: 7},
		{source: "Foo\nBar\n", bytes: 8},
		{source: "Foo\nBar\nBaz", bytes: 11},
	}

	for _, example := range examples {
		reader := strings.NewReader(example.source)
		result := count(reader)
		if result.bytes != example.bytes {
			t.Errorf("count(\"%s\") result.bytes => %d, want %d", example.source, result.bytes, example.bytes)
		}
	}
}

func TestWriteResults(t *testing.T) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	expected := "      10     256    2345\n    1011   25644234523232\n"
	var results []result
	results = append(results, result{lines: 10, words: 256, bytes: 2345})
	results = append(results, result{lines: 1011, words: 25644, bytes: 234523232})

	writeResults(writer, results)

	if output.String() != expected {
		t.Errorf("writeResults() => \"%s\", want \"%s\"", output.String(), expected)
	}
}
