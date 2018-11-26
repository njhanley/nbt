package nbt

import (
	"bytes"
	"testing"
)

func TestEncodeSorted(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.SortCompounds(true)
	if err := enc.Encode(testTag); err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()
	for i, n := range data {
		if i >= len(testData) || n != testData[i] {
			t.Fatalf("expected and got differ at byte %d:\nexpected: %#02v\ngot: %#02v\n", i, testData, data)
		}
	}
}
