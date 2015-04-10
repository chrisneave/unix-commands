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
	tailedFiles := make([]*tailedFile, flag.NArg())
	writer := bufio.NewWriter(os.Stdout)
	var err error
	var fi os.FileInfo

	for i, arg := range flag.Args() {
		tailedFiles[i] = &tailedFile{filename: arg}

		tailedFiles[i].file, err = os.Open(tailedFiles[i].filename)
		if err != nil {
			log.Fatal(err)
		}
		defer tailedFiles[i].file.Close()

		fi, err = tailedFiles[i].file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		tailedFiles[i].lastFileSize = fi.Size()

		if len(tailedFiles) > 1 {
			tailedFiles[i].writeHeaderTo(writer)
		}

		tailedFiles[i].offset = writeTail(tailedFiles[i].file, writer, *limit)
	}

	for {
		for _, tf := range tailedFiles {
			if tf.hasChanged() == false {
				continue
			}

			fi, err = tf.file.Stat()
			if err != nil {
				log.Fatal(err)
			}
			tf.lastFileSize = fi.Size()
			if len(tailedFiles) > 1 {
				tf.writeHeaderTo(writer)
			}

			reader := bufio.NewReader(tf.file)
			for {
				line, err := reader.ReadString(10)
				if len(line) > 0 {
					writer.WriteString(line)
				}
				if err != nil {
					writer.Flush()
					break
				}
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
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

// Write the last n lines from the input to the output, no seeking.
func writeTail(input io.Reader, output io.Writer, lineCount int) (newOffset int64) {
	lines, newOffset := tailScan(input, lineCount, 0)
	writer := bufio.NewWriter(output)
	for _, line := range lines {
		writer.WriteString(line)
		writer.WriteString("\n")
	}
	writer.Flush()
	return
}

func (tf *tailedFile) writeHeaderTo(writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("==> %s <==", tf.filename))
	writer.WriteString("\n")
}

func (tf *tailedFile) hasChanged() bool {
	fi, err := tf.file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fi.Size() == tf.lastFileSize {
		return false
	}

	return true
}
