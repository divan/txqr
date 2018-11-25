package txqr

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
		t.Fatalf("Encode failed: %v", err)
	}
	if len(chunks) != 2 {
		t.Fatalf("Expected 2 chunks")
	}

	expected := "0/101|" + strings.Repeat("s", N)
	if chunks[0] != expected {
		t.Fatalf("First chunk is invalid: '%s'", chunks[0])
	}
	expected = "100/101|" + "s"
	if chunks[1] != expected {
		t.Fatalf("Second chunk is invalid: '%s'", chunks[1])
	}
}
