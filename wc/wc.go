package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	var results []result

	for _, arg := range flag.Args() {
		file, err := os.Open(arg)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		results = append(results, count(file))
		writeResults(os.Stdout, results)
	}

	if len(results) == 0 {
		results = append(results, count(os.Stdin))
		writeResults(os.Stdout, results)
	}
}

type result struct {
	lines int64
	words int64
	bytes int64
}

func count(input io.Reader) (r result) {
	lineReader := bufio.NewReader(input)
	var newLine byte = 10

	for {
		line, err := lineReader.ReadString(newLine)
		if len(line) > 0 {
			r.lines++
		}
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
	for _, r := range results {
		writer.WriteString(fmt.Sprintf("%8d%8d%8d\n", r.lines, r.words, r.bytes))
	}
	writer.Flush()
}
