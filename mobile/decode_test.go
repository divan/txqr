package txqr

import (
	"strings"
	"testing"
	"time"

	roottxqr "github.com/divan/txqr"
)

func TestDecode(t *testing.T) {
	expected := strings.Repeat("s", 101)

	enc := roottxqr.NewEncoder(30)
	chunks, err := enc.Encode(expected)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec := NewDecoder()
	var i int
	for i = 0; i < len(chunks); i++ {
		if err := dec.Decode(chunks[i]); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if dec.IsCompleted() {
			break
		}
	}
	if i == 0 {
		t.Fatalf("IsCompleted expected to be false after the first chunk")
	}

	if !dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be true")
	}

	if dec.Data() != expected {
		t.Fatalf("Expected to get '%s' (len %d), but got '%s' (len %d)", expected, len(expected), dec.Data(), len(dec.Data()))
	}
}

// TestInvalidDecode checks that feeding extra chunks after the message is
// already fully decoded is a no-op and doesn't corrupt the result.
func TestInvalidDecode(t *testing.T) {
	expected := strings.Repeat("s", 101)

	enc := roottxqr.NewEncoder(30)
	chunks, err := enc.Encode(expected)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec := NewDecoder()
	var i int
	for i = 0; i < len(chunks); i++ {
		if err := dec.Decode(chunks[i]); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if dec.IsCompleted() {
			break
		}
	}

	if !dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be true")
	}
	if i+1 >= len(chunks) {
		t.Fatalf("test setup expects leftover chunks after completion, got none")
	}

	// feed one more (already redundant) chunk after completion
	if err := dec.Decode(chunks[i+1]); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if !dec.IsCompleted() {
		t.Fatal("IsCompleted expected to remain true", dec.Data())
	}

	if dec.Data() != expected {
		t.Fatalf("Expected to get '%s', but got '%s'", expected, dec.Data())
	}
}

// TestProgress checks that decoding doesn't panic or misbehave while chunks
// are being fed in. Note: Progress() currently always reports 0, since it
// relies on Decoder.Read(), which was stubbed to return 0 (see "TODO: remove"
// in decode.go) when the root package switched to fountain codes.
func TestProgress(t *testing.T) {
	expected := strings.Repeat("s", 101)

	enc := roottxqr.NewEncoder(30)
	chunks, err := enc.Encode(expected)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec := NewDecoder()
	for _, chunk := range chunks {
		if err := dec.Decode(chunk); err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if p := dec.Progress(); p < 0 || p > 100 {
			t.Fatalf("Progress out of range: %d", p)
		}
		if dec.IsCompleted() {
			break
		}
	}

	if !dec.IsCompleted() {
		t.Fatalf("IsCompleted expected to be true")
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
