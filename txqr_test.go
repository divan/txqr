package txqr

import (
	"fmt"
	"strings"
	"testing"
)

func TestTXQR(t *testing.T) {
	str := strings.Repeat("hello, world!", 1000)
	enc := NewEncoder(10)
	chunks, err := enc.Encode(str)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec := NewDecoder()
	for _, chunk := range chunks {
		err = dec.Decode(chunk)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
	}
	got := dec.Data()
	if got != str {
		t.Fatalf("Expected '%s', but got '%s'", str, got)
	}
}

func BenchmarkTXQREncode(b *testing.B) {
	var tests = []struct {
		length  int
		chunkSz int
	}{
		{100, 10},
		{1000, 10},
		{1000, 100},
		{10000, 100},
		{10000, 1000},
	}

	for _, test := range tests {
		b.Run(fmt.Sprintf("%d, %d", test.length, test.chunkSz), func(b *testing.B) {
			str := strings.Repeat("hello, world!", test.length)
			enc := NewEncoder(test.chunkSz)
			for i := 0; i < b.N; i++ {
				_, _ = enc.Encode(str)
			}
		})
	}
}

func BenchmarkTXQRDecode(b *testing.B) {
	var tests = []struct {
		length  int
		chunkSz int
	}{
		{100, 10},
		{1000, 10},
		{1000, 100},
		{10000, 100},
		{10000, 1000},
	}

	for _, test := range tests {
		b.Run(fmt.Sprintf("%d, %d", test.length, test.chunkSz), func(b *testing.B) {
			str := strings.Repeat("hello, world!", test.length)
			enc := NewEncoder(test.chunkSz)
			chunks, err := enc.Encode(str)
			if err != nil {
				b.Fatalf("Encode failed: %v", err)
			}
			dec := NewDecoder()
			for i := 0; i < b.N; i++ {
				for _, chunk := range chunks {
					err = dec.Decode(chunk)
					if err != nil {
						b.Fatalf("Decode failed: %v", err)
					}
				}
			}
		})
	}
}
