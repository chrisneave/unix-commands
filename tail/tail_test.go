package main

import "testing"

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
		t.Errorf("appendAndTail(%s, %s, %d) => length of result is %d, want %d", dest, src, limit, len(result), expectedLength)
		return
	}

	// Check content
	for i, v := range result {
		if expected[i] != v {
			t.Errorf("appendAndTail(%s, %s, %d) => got %s, want %s", dest, src, limit, v, expected[i])
		}
	}
}

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

func stringBuffersAreEqual(a, b []string) bool {
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
	if stringBuffersAreEqual(expected, lines) == false {
		t.Errorf("Expected %s but got %s", expected, lines)
	}
}

func TestTailScanWithLargeLimit(t *testing.T) {
	expected := []string{"I have", "three", "lines"}
	reader := strings.NewReader("I have\nthree\nlines")
	lines, _ := tailScan(reader, 10, 0)
	if stringBuffersAreEqual(expected, lines) == false {
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
