package squish

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	const input = "WWWWWWWWWWWWBWWWWWWWWWWWWBBBWWWWWWWWWWWWWWWWWWWWWWWWBWWWWWWWWWWWWWW"
	s := strings.NewReader(input)
	buf := &bytes.Buffer{}
	if err := Encode(buf, s); err != nil {
		t.Fatal(err)
	}
	if buf.Len() > len(input) {
		t.Errorf("encoded length should not be greater than %d bytes but it was %d", len(input), buf.Len())
	}
	result := &bytes.Buffer{}
	if err := Decode(buf, result); err != nil {
		t.Fatal(err)
	}
	if result.String() != input {
		t.Errorf("expecting: %s\nreceived: %s\n", input, result.String())
	}
}
