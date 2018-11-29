package txqr

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestTXQR(t *testing.T) {
	str := strings.Repeat("hello, world!", 1000)
	enc := NewEncoder(10)
	chunks, err := enc.Encode(str)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec := NewDecoder()
	for !dec.IsCompleted() {
		for _, chunk := range chunks {
			err = dec.Decode(chunk)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
		}
	}
	got := dec.Data()
	if got != str {
		t.Fatalf("Expected '%s', but got '%s'", str, got)
	}
}

// TestTXQRErasures tests information decoding over erasure channels
// with different erasure probabilities.
func TestTXQRErasures(t *testing.T) {
	str := strings.Repeat("hello, world!", 1000)
	enc := NewEncoder(10)
	chunks, err := enc.Encode(str)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	dec := NewDecoder()
	now := time.Now()
	for !dec.IsCompleted() {
		// erase new set of chunks every time
		toDel := chunksCountToDelete(len(chunks))
		transmitted := eraseChunks(chunks, toDel)
		for _, chunk := range transmitted {
			err = dec.Decode(chunk)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
		}
	}
	duration := time.Since(now)
	t.Logf("Sending over erasure channel took: %v", duration)
	got := dec.Data()
	if got != str {
		t.Fatalf("Expected '%s', but got '%s'", str, got)
	}
}

// eraseChunks randomly erases n chunks from the input slice.
func eraseChunks(chunks []string, n int) []string {
	toErase := make(map[int]bool)
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(chunks))
		toErase[idx] = true
	}

	ret := make([]string, 0, len(chunks)-n)
	for idx, chunk := range chunks {
		if toErase[idx] {
			continue
		}
		ret = append(ret, chunk)
	}
	return ret
}

// chunksCountToDelete returns random number of chunks for deletion,
// using normal distribution with 2 std deviation and N/3 as a mean.
func chunksCountToDelete(n int) int {
	mean := float64(n / 3)
	dev := 2.0
	del := int(rand.NormFloat64()*dev + mean)
	if del < 0 {
		del = 0
	}
	if del > n {
		del = n
	}
	return del
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

// BenchmarkTXQRErasures benchmarks information decoding over erasure channels
// with different erasure probabilities.
func BenchmarkTXQRErasures(b *testing.B) {
	var tests = []struct {
		length  int
		chunkSz int
	}{
		{10000, 10},
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

			for i := 0; i < b.N; i++ {
				dec := NewDecoder()
				for !dec.IsCompleted() {
					// erase new set of chunks every time
					toDel := chunksCountToDelete(len(chunks))
					transmitted := eraseChunks(chunks, toDel)
					for _, chunk := range transmitted {
						err = dec.Decode(chunk)
						if err != nil {
							b.Fatalf("Decode failed: %v", err)
						}
					}
				}
			}
		})
	}
}
