package txqr

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/divan/txqr/qr"
)

func TestTXQR(t *testing.T) {
	str := strings.Repeat("hello, world!", 1000)
	ch, err := Encode(str, 100, 512, qr.Medium)
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

func TestTXQRBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping duplex benchmark in non-short test")
		return
	}
	filename := "qr/testdata/15k.jpg"
	runDuplexBenchmark(t, filename)
}

func runDuplexBenchmark(t *testing.T, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Read test file: %v", err)
	}

	for _, lvl := range []qr.RecoveryLevel{qr.Low, qr.Medium, qr.High, qr.Highest} {
		for split := 50; split < 1000; split += 150 {
			dur := runSinglePass(t, data, split, 512, lvl)
			fmt.Printf("Level %s, Split %d: %v\n", lvlToStr(lvl), split, dur)
		}
	}
}

func runSinglePass(t *testing.T, data []byte, split, size int, lvl qr.RecoveryLevel) time.Duration {
	now := time.Now()
	str := base64.StdEncoding.EncodeToString(data)
	ch, err := Encode(string(str), split, size, lvl)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	gotB64, err := Decode(ch)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	got, err := base64.StdEncoding.DecodeString(gotB64)
	if err != nil {
		t.Fatalf("Base64 Decode failed: %v", err)
	}
	if len(got) != len(data) {
		t.Fatalf("Expected %d bytes, but got %d", len(data), len(got))
	}
	if string(got) != string(data) {
		t.Fatalf("Result doesn't match original")
	}
	return time.Since(now)
}

func lvlToStr(lvl qr.RecoveryLevel) string {
	switch lvl {
	case qr.Low:
		return "low"
	case qr.Medium:
		return "medium"
	case qr.High:
		return "high"
	case qr.Highest:
		return "highest"
	default:
		return "unknown"
	}
}
