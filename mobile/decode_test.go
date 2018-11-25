package txqr

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	N := 100
	expected := strings.Repeat("s", N+1)
	dec := NewDecoder()
	var err error
	err = dec.Decode("0/101|" + strings.Repeat("s", N))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be false")
	}
	err = dec.Decode("100/101|" + "s")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if dec.IsCompleted() == false {
		t.Fatal("IsCompleted expected to be true")
	}

	if dec.Data() != expected {
		t.Fatalf("Expected to get '%s' (len %d), but got '%s' (len %d)", expected, len(expected), dec.Data(), len(dec.Data()))
	}
}

func TestInvalidDecode(t *testing.T) {
	N := 100
	expected := strings.Repeat("s", N+1)
	dec := NewDecoder()
	var err error
	err = dec.Decode("0/101|" + strings.Repeat("s", 90))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be false")
	}
	if dec.Progress() != 89 {
		t.Fatalf("Progress should be equal to %v, but got %v", 99, dec.Progress())
	}
	err = dec.Decode("100/101|" + "s")
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be false")
	}
	// missing out-of-order chunk
	err = dec.Decode("90/101|" + strings.Repeat("s", 10))
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if dec.IsCompleted() == false {
		t.Fatal("IsCompleted expected to be true", dec.Data())
	}

	if dec.Data() != expected {
		t.Fatalf("Expected to get '%s', but got '%s'", expected, dec.Data())
	}
}

func TestProgress(t *testing.T) {
	N := 100
	dec := NewDecoder()
	var err error

	for i := 0; i < N; i++ {
		chunk := fmt.Sprintf("%d/%d|s", i, N)
		err = dec.Decode(chunk)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if dec.Progress() != i+1 {
			t.Fatalf("Progress should be equal to %v, but got %v", i+1, dec.Progress())
		}
	}
}

func TestTotalTime(t *testing.T) {
	dur := 12345678 * time.Microsecond // 12.345678s
	got := formatDuration(dur)
	expected := "12.3s"
	if got != expected {
		t.Fatalf("Expected str to be '%s', but got '%s'", expected, got)
	}
}
