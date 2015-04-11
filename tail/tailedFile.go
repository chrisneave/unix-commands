package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type file interface {
	Stat() (fi os.FileInfo, err error)
	io.Reader
	io.Closer
}

type tailedFile struct {
	filename     string
	file         file
	offset       int64
	lastFileSize int64
}

func newTailedFileFromFile(f file) *tailedFile {
	return &tailedFile{file: f}
}

func (tf *tailedFile) Stat() (fi os.FileInfo, err error) {
	return tf.file.Stat()
}

func (tf *tailedFile) writeHeaderTo(writer *bufio.Writer) {
	fi, _ := tf.file.Stat()
	writer.WriteString(fmt.Sprintf("==> %s <==", fi.Name()))
	writer.WriteString("\n")
	writer.Flush()
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
