package main

import (
	"bufio"
	"bytes"
	"testing"
)

func createSubject(file file) *tailedFile {
	tf := newTailedFileFromFile(file)
	if tf == nil {
		panic("TailedFile is nil")
	}
	return tf
}

func TestTailedFileStatReturnsFileInfo(t *testing.T) {
	ff := &FakeFile{name: "foo.txt"}
	tf := createSubject(ff)

	fi, _ := tf.Stat()
	if fi == nil {
		t.Error("FileInfo is nil")
	}
	if fi.Name() != ff.name {
		t.Errorf("Expected %s but got %s", ff.name, fi.Name())
	}
}

func TestTailedFileWriteHeaderTo(t *testing.T) {
	var output bytes.Buffer
	writer := bufio.NewWriter(&output)
	expected := "==> foo.txt <==\n"
	ff := &FakeFile{name: "foo.txt"}
	tf := createSubject(ff)

	tf.writeHeaderTo(writer)

	if output.String() != expected {
		t.Errorf("Expected %s but got %s", expected, output.String())
	}
}

func TestTailedFileHasNotChanged(t *testing.T) {
	ff := &FakeFile{name: "foo.txt"}
	tf := createSubject(ff)

	actual := tf.hasChanged()

	if actual != false {
		t.Error("File has changed when it should not")
	}
}

func TestTailedFileHasChanged(t *testing.T) {
	ff := &FakeFile{name: "foo.txt", size: 100}
	tf := createSubject(ff)
	tf.lastFileSize = 50

	actual := tf.hasChanged()

	if actual != true {
		t.Error("File has not changed when it should have")
	}
}

func TestWriteTailToWritesLinesToWriter(t *testing.T) {
	examples := []struct {
		source         string
		lines          int
		expectedOutput string
		expectedOffset int64
	}{
		{source: "", lines: 1, expectedOutput: "", expectedOffset: 0},
		{source: "Foo\n", lines: 1, expectedOutput: "Foo\n", expectedOffset: 4},
		{source: "Foo", lines: 1, expectedOutput: "Foo\n", expectedOffset: 3},
		{source: "Foo\nBar\n", lines: 1, expectedOutput: "Bar\n", expectedOffset: 8},
		{source: "Foo\nBar\n", lines: 2, expectedOutput: "Foo\nBar\n", expectedOffset: 8},
		{source: "Foo\nBar\n", lines: 3, expectedOutput: "Foo\nBar\n", expectedOffset: 8},
		{source: "Foo\nBar\n", lines: 0, expectedOutput: "", expectedOffset: 8},
	}

	for _, example := range examples {
		ff := FakeFile{content: example.source}
		tf := newTailedFileFromFile(&ff)
		var output bytes.Buffer

		tf.writeTailTo(&output, example.lines)

		if example.expectedOutput != output.String() {
			t.Errorf("writeTailTo('%s', %d) => output = %s, want %s", example.source, example.lines, output.String(), example.expectedOutput)
		}

		if example.expectedOffset != tf.offset {
			t.Errorf("writeTailTo('%s', %d) => tf.offset = %d, want %d", example.source, example.lines, tf.offset, example.expectedOffset)
		}
	}
}
