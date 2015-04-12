package main

import (
	"bufio"
	"flag"
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

func main() {
	flag.Parse()
	tailedFiles := make([]*tailedFile, flag.NArg())
	writer := bufio.NewWriter(os.Stdout)
	var err error
	var fi os.FileInfo

	for i, arg := range flag.Args() {
		tf := &tailedFile{}
		tailedFiles[i] = tf

		tf.file, err = os.Open(arg)
		if err != nil {
			log.Fatal(err)
		}
		defer tf.file.Close()

		fi, err = tf.Stat()
		if err != nil {
			log.Fatal(err)
		}
		tf.lastFileSize = fi.Size()

		if len(tailedFiles) > 1 {
			tf.writeHeaderTo(writer)
		}

		tf.writeTailTo(writer, *limit)
	}

	if *follow == false {
		os.Exit(0)
	}

	for {
		for _, tf := range tailedFiles {
			if tf.hasChanged() == false {
				continue
			}

			fi, err = tf.Stat()
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
