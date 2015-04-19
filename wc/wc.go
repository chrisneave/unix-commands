package main

import (
	"bufio"
	"io"
	"strings"
)

type result struct {
	lines int64
	words int64
}

func count(input io.Reader) (r result) {
	lineReader := bufio.NewReader(input)
	var newLine byte = 10

	for {
		line, err := lineReader.ReadString(newLine)
		r.lines++

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
