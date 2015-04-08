package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	follow = flag.Bool("f", false, "Follow the files being tailed.")
	limit  = flag.Int("lines", 10, "The number of lines to output from the tailed file")
)

type tailedFile struct {
	filename     string
	file         *os.File
	offset       int64
	lastFileSize int64
}

func main() {
	flag.Parse()
	tailedFiles := make([]tailedFile, flag.NArg())
	var err error
	var header string

	for i, arg := range flag.Args() {
		tailedFiles[i] = tailedFile{filename: arg}

		tailedFiles[i].file, err = os.Open(tailedFiles[i].filename)
		if err != nil {
			log.Fatal(err)
		}
		defer tailedFiles[i].file.Close()

		if len(tailedFiles) > 1 {
			header = fmt.Sprintf("==> %s <==", tailedFiles[i].filename)
		}

		tailedFiles[i].offset = seekAndOutput(tailedFiles[i].file, tailedFiles[i].offset, header)
	}

	for {
		for _, tf := range tailedFiles {
			fs, err := os.Stat(tf.filename)
			if err != nil {
				log.Fatal(err)
			}

			tf.follow(fs.Size(), len(tailedFiles) > 1)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (tf *tailedFile) follow(currentSize int64, outputHeader bool) {
	if tf.lastFileSize < currentSize {
		tf.offset = seekAndOutput(tf.file, tf.offset, fmt.Sprintf("==> %s <==", tf.filename))
		tf.lastFileSize = currentSize
	}
}

func seekAndOutput(source io.ReadSeeker, offset int64, header string) (newOffset int64) {
	var currentLines []string

	offset, err := source.Seek(offset, 0)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(source)
	currentLines, newOffset = tailScan(reader, *limit, offset)

	if len(header) > 0 && len(currentLines) > 0 {
		fmt.Println(header)
	}

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
	reader := bufio.NewReader(source)
	for {
		line, err := reader.ReadString(10)
		line = strings.TrimSuffix(line, "\n")
		if len(line) > 0 {
			lines = appendAndTail(lines, []string{line}, limit)
		}
		newOffset += int64(len(line) + 1)
		if err != nil {
			break
		}
	}
	// Adjust the offset but ignore any trailing new line
	newOffset--
	return
}
