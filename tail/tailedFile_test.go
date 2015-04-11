package main

import (
	"io"
	"testing"
)

type StringFile struct {
	content string
	offset  int64
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
