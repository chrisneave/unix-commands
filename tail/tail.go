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

			if len(tailedFiles) > 1 {
				header = fmt.Sprintf("==> %s <==", tf.filename)
			}

			if tf.lastFileSize < fs.Size() {
				tf.offset = seekAndOutput(tf.file, tf.offset, header)
				tf.lastFileSize = fs.Size()
			}
		}

		time.Sleep(500 * time.Millisecond)
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

	if len(header) > 0 {
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
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = appendAndTail(lines, []string{scanner.Text()}, limit)
		newOffset += int64(len(scanner.Text()) + 1)
	}
	// Adjust the offset but ignore any trailing new line
	newOffset--
	return
}
