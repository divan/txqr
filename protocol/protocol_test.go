package protocol

import (
	"strings"
	"testing"
)

func TestEncodeString(t *testing.T) {
	N := 100
	enc := NewEncoder(N)
	str := strings.Repeat("s", N+1)
	chunks, err := enc.Encode(str)
	if err != nil {
		t.Fatalf("encodeString failed: %v", err)
	}
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks")
	}

	expected := "0/65|" + strings.Repeat("s", N)
	if chunks[0] != expected {
		t.Fatalf("First chunk is invalid: '%s'", chunks[0])
	}
	expected = "64/65|" + "s"
	if chunks[1] != expected {
		t.Fatalf("Second chunk is invalid: '%s'", chunks[1])
	}
}
