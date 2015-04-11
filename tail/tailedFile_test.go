package main

import "testing"

func TestTailedFileStatReturnsFileInfo(t *testing.T) {
	ff := &FakeFile{name: "foo.txt"}
	tf := newTailedFileFromFile(ff)
	if tf == nil {
		t.Fatal("TailedFile is nil")
	}
	fi, _ := tf.Stat()
	if fi == nil {
		t.Error("FileInfo is nil")
	}
	if fi.Name() != ff.name {
		t.Errorf("Expected %s but got %s", ff.name, fi.Name())
	}
}
