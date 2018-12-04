package nbt

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecoder(t *testing.T) {
	tag, err := NewDecoder(bytes.NewReader(testData)).Decode()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(testTag, tag); diff != "" {
		t.Fatalf("cmp.Diff(expected, got):\n%v", diff)
	}
}
