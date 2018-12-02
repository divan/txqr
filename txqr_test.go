package txqr

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestTXQR(t *testing.T) {
	var tests = []struct {
		length  int
		chunkSz int
	}{
		{10 * 1024, 100},
		{10 * 1024, 200},
		{10 * 1024, 300},
		{10 * 1024, 400},
		{10 * 1024, 500},
		{10 * 1024, 650},
		{10 * 1024, 800},
		{10 * 1024, 900},
		{10 * 1024, 1000},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d, %d", test.length, test.chunkSz), func(t *testing.T) {
			str := newTestData(test.length)
			enc := NewEncoder(test.chunkSz)
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
		})
	}
}

// TestTXQRErasures tests information decoding over erasure channels
// with different erasure probabilities.
func TestTXQRErasures(t *testing.T) {
	var tests = []struct {
		length  int
		chunkSz int
		fps     int
	}{
		{10 * 1024, 100, 3},
		{10 * 1024, 300, 3},
		{10 * 1024, 500, 3},
		{10 * 1024, 650, 3},
		{10 * 1024, 800, 3},
		{10 * 1024, 800, 3},
		{10 * 1024, 100, 6},
		{10 * 1024, 300, 6},
		{10 * 1024, 500, 6},
		{10 * 1024, 800, 6},
		{10 * 1024, 800, 6},
		{10 * 1024, 100, 9},
		{10 * 1024, 300, 9},
		{10 * 1024, 500, 9},
		{10 * 1024, 800, 9},
		{10 * 1024, 800, 9},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d, %d", test.length, test.chunkSz), func(t *testing.T) {
			str := newTestData(test.length)
			enc := NewEncoder(test.chunkSz)
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
					if dec.IsCompleted() {
						break
					}
				}
				time.Sleep(1 * time.Second / time.Duration(test.fps))
			}
			duration := time.Since(now)
			t.Logf("[%db/%d, %dfps] took: %v", test.length, test.chunkSz, test.fps, duration)
			got := dec.Data()
			if got != str {
				t.Fatalf("Expected '%s', but got '%s'", str, got)
			}
		})
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
	mean := float64(n / 2)
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
		{10 * 1024, 100},
		{10 * 1024, 1000},
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
		{10 * 1024, 100},
		{10 * 1024, 1000},
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
		{10 * 1024, 10},
		{10 * 1024, 100},
		{10 * 1024, 1000},
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

func newTestData(size int) string {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(fmt.Sprintf("Can't generate rand data: %v", err))
	}
	return string(data)
}
