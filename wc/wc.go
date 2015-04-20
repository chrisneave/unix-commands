package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type result struct {
	lines int64
	words int64
	bytes int64
}

func count(input io.Reader) (r result) {
	lineReader := bufio.NewReader(input)
	var newLine byte = 10

	// Use Peek to check if there is at least one character in
	// the input. If so then increment the byte count by one
	// to account for the trailing byte that the code below
	// will miss.
	_, err := lineReader.Peek(1)
	if err == nil {
		r.bytes++
	}

	for {
		line, err := lineReader.ReadString(newLine)
		r.lines++
		r.bytes += int64(len(line))

		wordScanner := bufio.NewScanner(strings.NewReader(line))
		wordScanner.Split(bufio.ScanWords)
		for wordScanner.Scan() {
			r.words++
		}

		if err != nil {
			break
		}
	}

	return
}

func writeResults(output io.Writer, results []result) {
	writer := bufio.NewWriter(output)
	for i, r := range results {
		writer.WriteString(fmt.Sprintf("%8d%8d%8d", r.lines, r.words, r.bytes))
		if i < len(results)-1 {
			writer.WriteString("\n")
		}
	}
	writer.Flush()
}
