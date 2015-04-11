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
}

type tailedFile struct {
	filename     string
	file         *os.File
	offset       int64
	lastFileSize int64
}

func (tf *tailedFile) Stat() (fi os.FileInfo, err error) {
	return tf.file.Stat()
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
