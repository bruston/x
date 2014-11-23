package main

import (
	"bytes"
	"testing"
)

var lines = []byte(`test line 1
a longer test line
another line in the test file
just keep testing`)

func TestLongest(t *testing.T) {
	r := bytes.NewBuffer(lines)
	expected := Result{
		text:   "another line in the test file",
		line:   3,
		length: 29,
	}
	result := longest(r)
	if result != expected {
		t.Errorf("expecting %#v, got %#v", expected, result)
	}
}
