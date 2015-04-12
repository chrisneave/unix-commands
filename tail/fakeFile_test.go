package main

import (
	"io"
	"os"
	"testing"
	"time"
)

type StringFile struct {
	content string
	offset  int64
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

type FakeFile StringFile

func (fr *FakeFile) Read(b []byte) (n int, err error) {
	buffer := []byte(fr.content)
	buffer = buffer[fr.offset:]
	n = copy(b, buffer)
	fr.offset += int64(n)
	if fr.offset >= int64(len(fr.content)) {
		return n, io.EOF
	}
	return n, nil
}

func (fr *FakeFile) Stat() (fi os.FileInfo, err error) {
	return fr, nil
}

func (fr *FakeFile) Name() string {
	return fr.name
}

func (fr *FakeFile) Size() int64 {
	return fr.size
}

func (fr *FakeFile) Mode() os.FileMode {
	return fr.mode
}

func (fr *FakeFile) ModTime() time.Time {
	return fr.modTime
}

func (fr *FakeFile) IsDir() bool {
	return fr.isDir
}

func (fr *FakeFile) Sys() interface{} {
	return nil
}

func (fr *FakeFile) Close() error {
	return nil
}

type FakeFileWithStatError StringFile

func (fr *FakeFileWithStatError) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (fr *FakeFileWithStatError) Stat() (fi os.FileInfo, err error) {
	return nil, &os.PathError{}
}

func (fr *FakeFileWithStatError) Close() error {
	return nil
}

func TestFakeFileReadReturnsFileContent(t *testing.T) {
	subject := FakeFile{content: "Some file data"}
	buffer := make([]byte, len(subject.content))

	bytesRead, err := subject.Read(buffer)

	if bytesRead != len(subject.content) {
		t.Errorf("Expected %d but got %d", len(subject.content), bytesRead)
	}

	if string(buffer) != subject.content {
		t.Errorf("Expected '%s' but got '%s'", string(buffer), subject.content)
	}

	if err != io.EOF {
		t.Errorf("Expected %s but got %s", io.EOF, err)
	}
}

func TestFakeFileReadStartsFromOffset(t *testing.T) {
	subject := FakeFile{content: "Some file data", offset: 5}
	buffer := make([]byte, len(subject.content))

	bytesRead, _ := subject.Read(buffer)

	if bytesRead != len(subject.content[5:]) {
		t.Errorf("Expected %d but got %d", len(subject.content[5:]), bytesRead)
	}

	if string(buffer[:bytesRead]) != subject.content[5:] {
		t.Errorf("Expected '%s' but got '%s'", string(buffer), subject.content[5:])
	}
}
