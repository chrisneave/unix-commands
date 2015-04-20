package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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

		r := count(file)
		r.filename = arg
		results = append(results, r)
	}

	if len(results) == 0 {
		results = append(results, count(os.Stdin))
	}

	writeResults(os.Stdout, results)
}

type result struct {
	lines    int64
	words    int64
	bytes    int64
	filename string
}

func count(input io.Reader) (r result) {
	lineReader := bufio.NewReader(input)
	newLine := []byte("\n")

	for {
		line, err := lineReader.ReadBytes(newLine[0])
		if len(line) > 0 {
			r.lines++
		}
		r.bytes += int64(len(line))

		wordScanner := bufio.NewScanner(bytes.NewReader(line))
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
	var total result
	writer := bufio.NewWriter(output)
	for _, r := range results {
		writer.WriteString(fmt.Sprintf("%8d%8d%8d", r.lines, r.words, r.bytes))
		if r.filename != "" {
			writer.WriteString(fmt.Sprintf(" %s", r.filename))
		}
		writer.WriteString("\n")
		total.lines += r.lines
		total.words += r.words
		total.bytes += r.bytes
	}
	if len(results) > 1 {
		writer.WriteString(fmt.Sprintf("%8d%8d%8d total\n", total.lines, total.words, total.bytes))
	}
	writer.Flush()
}
