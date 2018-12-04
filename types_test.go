package nbt

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarshalJSON(t *testing.T) {
	data, err := json.MarshalIndent(testTag, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(testJSON, data); diff != "" {
		t.Fatalf("cmp.Diff(expected, got):\n%v", diff)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tag := new(NamedTag)
	if err := json.Unmarshal(testJSON, tag); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(testTag, tag); diff != "" {
		t.Fatalf("cmp.Diff(expected, got):\n%v", diff)
	}
}
