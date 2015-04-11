package main

import (
	"bytes"
	"io"
	"testing"
)
import "strings"

type Example struct {
	limit    int
	expected []string
}

func testAppendAndTail(t *testing.T, dest []string, src []string, limit int, expected []string) {
	expectedLength := len(dest) + len(src)
	if limit < expectedLength {
		expectedLength = limit
	}

	result := appendAndTail(dest, src, limit)

	// Check length
	if len(result) != expectedLength {
		t.Errorf("Length of returned slice should be %d, got %d", expectedLength, len(result))
		return
	}

	// Check content
	for i, v := range result {
		if expected[i] != v {
			t.Errorf("Expected is %s, Limit is %d, %s should equal %s", expected, limit, v, expected[i])
		}
	}
}

// TestAppendAndLimit ...
func TestAppendAndTail(t *testing.T) {
	examples := []Example{
		Example{0, []string{}},
		Example{1, []string{"5"}},
		Example{2, []string{"4", "5"}},
		Example{3, []string{"3", "4", "5"}},
		Example{4, []string{"2", "3", "4", "5"}},
		Example{5, []string{"1", "2", "3", "4", "5"}},
		Example{6, []string{"0", "1", "2", "3", "4", "5"}}}
	dest := []string{"0", "1", "2"}
	src := []string{"3", "4", "5"}

	for _, example := range examples {
		testAppendAndTail(t, dest, src, example.limit, example.expected)
	}
}

func TestAppendAndTailEmptySrc(t *testing.T) {
	examples := []Example{
		Example{0, []string{}},
		Example{1, []string{"5"}},
		Example{2, []string{"4", "5"}},
		Example{3, []string{"3", "4", "5"}},
		Example{4, []string{"3", "4", "5"}},
		Example{5, []string{"3", "4", "5"}},
		Example{6, []string{"3", "4", "5"}}}
	dest := []string{}
	src := []string{"3", "4", "5"}

	for _, example := range examples {
		testAppendAndTail(t, dest, src, example.limit, example.expected)
	}
}

func TestAppendAndTailEmptyDest(t *testing.T) {
	examples := []Example{
		Example{0, []string{}},
		Example{1, []string{"2"}},
		Example{2, []string{"1", "2"}},
		Example{3, []string{"0", "1", "2"}},
		Example{4, []string{"0", "1", "2"}},
		Example{5, []string{"0", "1", "2"}},
		Example{6, []string{"0", "1", "2"}}}
	dest := []string{"0", "1", "2"}
	src := []string{}

	for _, example := range examples {
		testAppendAndTail(t, dest, src, example.limit, example.expected)
	}
}

func BenchmarkAppendAndTail(b *testing.B) {
	dest := []string{"0", "1", "2"}
	src := []string{"3", "4"}

	for i := 0; i < b.N; i++ {
		appendAndTail(dest, src, 5)
	}
}

func buffersAreEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for n, v := range a {
		if v != b[n] {
			return false
		}
	}

	return true
}

func TestTailScan(t *testing.T) {
	expected := []string{"lines"}
	reader := strings.NewReader("I have\nthree\nlines")
	lines, _ := tailScan(reader, 1, 0)
	if buffersAreEqual(expected, lines) == false {
		t.Errorf("Expected %s but got %s", expected, lines)
	}
}

func TestTailScanWithLargeLimit(t *testing.T) {
	expected := []string{"I have", "three", "lines"}
	reader := strings.NewReader("I have\nthree\nlines")
	lines, _ := tailScan(reader, 10, 0)
	if buffersAreEqual(expected, lines) == false {
		t.Errorf("Expected %s but got %s", expected, lines)
	}
}

func TestTailScanWithOffset(t *testing.T) {
	var oldOffset int64 = 10
	source := "I have\nthree\nlines"
	var expected = oldOffset + int64(len(source))
	reader := strings.NewReader(source)
	_, offset := tailScan(reader, 10, oldOffset)
	if offset != expected {
		t.Errorf("Expected %d but got %d", expected, offset)
	}
}

func TestTailScanWithOffsetAndTrailingNewLine(t *testing.T) {
	var oldOffset int64 = 10
	source := "I have\nthree\nlines\n"
	var expected = oldOffset + int64(len(source))
	reader := strings.NewReader(source)
	_, offset := tailScan(reader, 10, oldOffset)
	if offset != expected {
		t.Errorf("Expected %d but got %d", expected, offset)
	}
}

type writeTailExample struct {
	source   string
	lines    int
	expected string
}

func TestWriteTail(t *testing.T) {
	examples := []writeTailExample{
		writeTailExample{source: "", lines: 1, expected: ""},
		writeTailExample{source: "Foo\n", lines: 1, expected: "Foo\n"},
		writeTailExample{source: "Foo", lines: 1, expected: "Foo\n"},
		writeTailExample{source: "Foo\nBar\n", lines: 1, expected: "Bar\n"},
		writeTailExample{source: "Foo\nBar\n", lines: 2, expected: "Foo\nBar\n"},
		writeTailExample{source: "Foo\nBar\n", lines: 3, expected: "Foo\nBar\n"},
		writeTailExample{source: "Foo\nBar\n", lines: 0, expected: ""}}

	for _, example := range examples {
		var output bytes.Buffer
		reader := strings.NewReader(example.source)
		writeTail(reader, &output, example.lines)

		if example.expected != output.String() {
			t.Errorf("Expected %s but got %s", example.expected, output.String())
		}
	}
}

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
