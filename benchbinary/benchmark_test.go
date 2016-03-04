package main

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type example struct {
	Uint8  uint8
	Uint32 uint32
}

func handleErr(b *testing.B, err error) {
	if err != nil {
		b.Fatalf("write failed: %s", err)
	}
}

func BenchmarkStruct(b *testing.B) {
	buf := bytes.Buffer{}
	e := example{
		10,
		10243,
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf := bytes.Buffer{}
		if err := binary.Write(&buf, binary.LittleEndian, e); err != nil {
			handleErr(b, err)
		}
	}
}

func BenchmarkIndividualFields(b *testing.B) {
	buf := bytes.Buffer{}
	e := example{
		10,
		10243,
	}
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := binary.Write(&buf, binary.LittleEndian, e.Uint8); err != nil {
			handleErr(b, err)
		}
		if err := binary.Write(&buf, binary.LittleEndian, e.Uint32); err != nil {
			handleErr(b, err)
		}
	}
}
