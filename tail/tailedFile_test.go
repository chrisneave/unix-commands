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
