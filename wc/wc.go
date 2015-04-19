package main

import (
	"bufio"
	"io"
)

type result struct {
	lines int64
}

func count(input io.Reader) (r result) {
	reader := bufio.NewReader(input)
	var newLine byte = 10

	for {
		_, err := reader.ReadString(newLine)
		r.lines++
		if err != nil {
			break
		}
	}

	return
}
