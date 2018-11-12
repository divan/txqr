package txqr

import (
	"strings"
	"testing"
)

func TestTXQR(t *testing.T) {
	str := strings.Repeat("hello, world!", 1000)
	ch, err := Encode(str)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	got, err := Decode(ch)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if got != str {
		t.Fatalf("Expected '%s', but got '%s'", str, got)
	}
}
