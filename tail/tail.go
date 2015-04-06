package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	follow = flag.Bool("f", false, "Follow the files being tailed.")
	limit  = flag.Int("lines", 10, "The number of lines to output from the tailed file")
)

func main() {
	flag.Parse()
	var filename = flag.Arg(0)
	var offset int64
	var lastFileSize int64

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	offset = seekAndOutput(file, offset)

	for {
		fs, err := os.Stat(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		if lastFileSize < fs.Size() {
			offset = seekAndOutput(file, offset)
			lastFileSize = fs.Size()
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func seekAndOutput(source io.ReadSeeker, offset int64) (newOffset int64) {
	var currentLines []string

	offset, err := source.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(source)
	currentLines, newOffset = tailScan(reader, *limit, offset)

	for _, line := range currentLines {
		fmt.Println(line)
	}

	return
}

func appendAndTail(dst []string, src []string, limit int) []string {
	combined := append(dst, src...)
	// Limit what we return based on the length of the combined slice.
	if limit > len(combined) {
		limit = len(combined)
	}
	// Tail the combined slice.
	return combined[len(combined)-limit:]
}

func tailScan(source io.Reader, limit int, offset int64) (lines []string, newOffset int64) {
	newOffset = offset
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = appendAndTail(lines, []string{scanner.Text()}, limit)
		newOffset += int64(len(scanner.Text()) + 1)
	}
	// Adjust the offset but ignore any trailing new line
	newOffset--
	return
}
