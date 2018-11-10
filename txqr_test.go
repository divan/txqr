package txqr

import "testing"

func TestTXQR(t *testing.T) {
	str := "hello, world!"
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
