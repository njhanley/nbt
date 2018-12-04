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
	if diff := cmp.Diff(testData, data); diff != "" {
		t.Fatalf("cmp.Diff(expected, got):\n%v", diff)
	}
}
