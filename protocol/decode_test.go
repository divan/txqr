package protocol

import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	N := 100
	expected := strings.Repeat("s", N+1)
	dec := NewDecoder()
	var err error
	err = dec.DecodeChunk("0/65|" + strings.Repeat("s", N))
	if err != nil {
		t.Fatalf("DecodeChunk failed: %v", err)
	}
	if dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be false")
	}
	err = dec.DecodeChunk("64/65|" + "s")
	if err != nil {
		t.Fatalf("DecodeChunk failed: %v", err)
	}
	if dec.IsCompleted() == false {
		t.Fatal("IsCompleted expected to be true")
	}

	if dec.Data() != expected {
		t.Fatalf("Expected to get '%s', but got '%s'", expected, dec.Data())
	}
}

func TestInvalidDecode(t *testing.T) {
	N := 100
	expected := strings.Repeat("s", N+1)
	dec := NewDecoder()
	var err error
	err = dec.DecodeChunk("0/65|" + strings.Repeat("s", 90))
	if err != nil {
		t.Fatalf("DecodeChunk failed: %v", err)
	}
	if dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be false")
	}
	err = dec.DecodeChunk("64/65|" + "s")
	if err != nil {
		t.Fatalf("DecodeChunk failed: %v", err)
	}
	if dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be false")
	}
	// missing out-of-order chunk
	err = dec.DecodeChunk("5a/65|" + strings.Repeat("s", 10))
	if err != nil {
		t.Fatalf("DecodeChunk failed: %v", err)
	}
	if dec.IsCompleted() == false {
		t.Fatal("IsCompleted expected to be true", dec.Data())
	}

	if dec.Data() != expected {
		t.Fatalf("Expected to get '%s', but got '%s'", expected, dec.Data())
	}
}
