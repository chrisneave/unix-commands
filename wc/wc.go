package main

import (
	"bufio"
	"io"
)

type result struct {
	lines int64
}

func count(input io.Reader) (r result) {
	lineScanner := bufio.NewScanner(input)
	for lineScanner.Scan() {
		r.lines++
	}
	return
}
