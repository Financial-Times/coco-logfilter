package main

import (
	"testing"
)

func TestFixBytesToString(t *testing.T) {
	// happy path
	input := []interface{}{float64('A'), float64('B')}
	output := fixBytesToString(input)
	expected := "AB"
	if output != expected {
		t.Errorf("expected output %v but got %v\n", expected, output)
	}
}

const expectedNewlines = `A
B
C`

func TestFixNewlines(t *testing.T) {
	input := "A|B|C"
	output := fixNewLines(input)
	if output != expectedNewlines {
		t.Errorf("expected %v but got %v\n", expectedNewlines, output)
	}
}
